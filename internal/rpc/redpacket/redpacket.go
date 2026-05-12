package redpacket

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/redpacket/chain"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	pbredpacket "github.com/openimsdk/protocol/redpacket"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"google.golang.org/grpc"
)

type Config struct {
	RpcConfig     config.RedPacket
	MongodbConfig config.Mongo
	Share         config.Share
	Discovery     config.Discovery
}

type redPacketServer struct {
	pbredpacket.UnimplementedRedPacketServer
	config         *Config
	db             controller.RedPacketDatabase
	chainClient    *chain.ChainClient
	tronClient     *chain.TronClient
	signerKey      *ecdsa.PrivateKey
	groupClient    *rpcli.GroupClient
	relationClient *rpcli.RelationClient
}

func Start(ctx context.Context, conf *Config, registry discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgoClient, err := mongoutil.NewMongoDB(ctx, conf.MongodbConfig.Build())
	if err != nil {
		return err
	}
	db := mgoClient.GetDB()

	rpDB, err := mgo.NewRedPacketMongo(db)
	if err != nil {
		return err
	}
	claimDB, err := mgo.NewRedPacketClaimMongo(db)
	if err != nil {
		return err
	}
	claimAuthDB, err := mgo.NewRedPacketClaimAuthMongo(db)
	if err != nil {
		return err
	}
	refundDB, err := mgo.NewRedPacketRefundMongo(db)
	if err != nil {
		return err
	}
	challengeDB, err := mgo.NewWalletBindingChallengeMongo(db)
	if err != nil {
		return err
	}
	bindingDB, err := mgo.NewWalletBindingMongo(db)
	if err != nil {
		return err
	}
	auditLogDB, err := mgo.NewAdminAuditLogMongo(db)
	if err != nil {
		return err
	}

	repo := controller.NewRedPacketDatabase(rpDB, claimDB, claimAuthDB, refundDB, challengeDB, bindingDB, auditLogDB)

	chainClient, err := chain.NewClient(
		conf.RpcConfig.Chain.RPCURL,
		conf.RpcConfig.Chain.ContractAddress,
		conf.RpcConfig.Chain.ChainID,
		conf.RpcConfig.Chain.SignerPrivateKey,
		conf.RpcConfig.Chain.ConfigAdminPrivateKey,
	)
	if err != nil {
		log.ZWarn(ctx, "redpacket eth client init failed, continuing without it", err)
		chainClient = nil
	}

	var tronClient *chain.TronClient
	if conf.RpcConfig.Tron.FullNodeURL != "" {
		abiJSON, abiErr := chain.ExtractABIFromEmbeddedArtifact()
		if abiErr != nil {
			log.ZWarn(ctx, "redpacket tron load abi failed", abiErr)
		} else {
			tronClient, err = chain.NewTronClient(
				conf.RpcConfig.Tron.FullNodeURL,
				conf.RpcConfig.Tron.ContractBase58,
				conf.RpcConfig.Tron.OwnerBase58,
				conf.RpcConfig.Tron.PrivateKeyHex,
				abiJSON,
				conf.RpcConfig.Tron.FeeLimit,
			)
			if err != nil {
				log.ZWarn(ctx, "redpacket tron client init failed", err)
				tronClient = nil
			}
		}
	}

	var signerKey *ecdsa.PrivateKey
	if k := conf.RpcConfig.Chain.SignerPrivateKey; k != "" {
		sk, parseErr := crypto.HexToECDSA(k)
		if parseErr != nil {
			log.ZWarn(ctx, "redpacket signer private key parse failed", parseErr)
		} else {
			signerKey = sk
		}
	}

	groupConn, err := registry.GetConn(ctx, conf.Share.RpcRegisterName.Group)
	if err != nil {
		return err
	}
	friendConn, err := registry.GetConn(ctx, conf.Share.RpcRegisterName.Friend)
	if err != nil {
		return err
	}

	srv := &redPacketServer{
		config:         conf,
		db:             repo,
		chainClient:    chainClient,
		tronClient:     tronClient,
		signerKey:      signerKey,
		groupClient:    rpcli.NewGroupClient(groupConn),
		relationClient: rpcli.NewRelationClient(friendConn),
	}

	pbredpacket.RegisterRedPacketServer(server, srv)

	if chainClient != nil {
		ethIndexer := chain.NewIndexer(chainClient, repo, conf.RpcConfig.Indexer.PollInterval, 0)
		ethIndexer.Start(ctx)
	}
	if tronClient != nil {
		tronIndexer := chain.NewTronIndexer(tronClient, repo, conf.RpcConfig.Indexer.PollInterval, 0)
		tronIndexer.Start(ctx)
	}

	return nil
}
