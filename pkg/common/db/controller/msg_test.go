package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_BatchInsertChat2DB(t *testing.T) {
	config.Config.Mongo.DBAddress = []string{"192.168.44.128:37017"}
	config.Config.Mongo.DBTimeout = 60
	config.Config.Mongo.DBDatabase = "openIM"
	config.Config.Mongo.DBSource = "admin"
	config.Config.Mongo.DBUserName = "root"
	config.Config.Mongo.DBPassword = "openIM123"
	config.Config.Mongo.DBMaxPoolSize = 100
	config.Config.Mongo.DBRetainChatRecords = 3650
	config.Config.Mongo.ChatRecordsClearTime = "0 2 * * 3"

	mongo, err := unrelation.NewMongo()
	if err != nil {
		t.Fatal(err)
	}
	err = mongo.GetDatabase().Client().Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	db := &commonMsgDatabase{
		msgDocDatabase: unrelation.NewMsgMongoDriver(mongo.GetDatabase()),
	}

	//ctx := context.Background()
	//msgs := make([]*sdkws.MsgData, 0, 1)
	//for i := 0; i < cap(msgs); i++ {
	//	msgs = append(msgs, &sdkws.MsgData{
	//		Content:  []byte(fmt.Sprintf("test-%d", i)),
	//		SendTime: time.Now().UnixMilli(),
	//	})
	//}
	//err = db.BatchInsertChat2DB(ctx, "test", msgs, 0)
	//if err != nil {
	//	panic(err)
	//}

	_ = db.BatchInsertChat2DB
	c := mongo.GetDatabase().Collection("msg")

	ch := make(chan int)
	rand.Seed(time.Now().UnixNano())

	index := 10

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(channelID int) {
			defer wg.Done()
			<-ch
			var arr []string
			for i := 0; i < 500; i++ {
				arr = append(arr, strconv.Itoa(i+1))
			}
			rand.Shuffle(len(arr), func(i, j int) {
				arr[i], arr[j] = arr[j], arr[i]
			})
			for j, s := range arr {
				if j == 0 {
					fmt.Printf("channnelID: %d, arr[0]: %s\n", channelID, arr[j])
				}
				filter := bson.M{"doc_id": "test:0"}
				update := bson.M{
					"$addToSet": bson.M{
						fmt.Sprintf("msgs.%d.del_list", index): bson.M{"$each": []string{s}},
					},
				}
				_, err := c.UpdateOne(context.Background(), filter, update)
				if err != nil {
					t.Fatal(err)
				}
			}
		}(i)
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ch
			var arr []string
			for i := 0; i < 500; i++ {
				arr = append(arr, strconv.Itoa(1001+i))
			}
			rand.Shuffle(len(arr), func(i, j int) {
				arr[i], arr[j] = arr[j], arr[i]
			})
			for _, s := range arr {
				filter := bson.M{"doc_id": "test:0"}
				update := bson.M{
					"$addToSet": bson.M{
						fmt.Sprintf("msgs.%d.read_list", index): bson.M{"$each": []string{s}},
					},
				}
				_, err := c.UpdateOne(context.Background(), filter, update)
				if err != nil {
					t.Fatal(err)
				}
			}
		}()
	}

	time.Sleep(time.Second * 2)

	close(ch)

	wg.Wait()

}

func TestName(t *testing.T) {
	s := ` [
        189,
        498,
        310,
        163,
        313,
        335,
        327,
        342,
        123,
        97,
        4,
        362,
        210,
        298,
        436,
        9,
        369,
        432,
        132,
        69,
        248,
        93,
        91,
        112,
        145,
        194,
        84,
        443,
        179,
        241,
        257,
        237,
        169,
        460,
        33,
        441,
        126,
        187,
        390,
        402,
        51,
        35,
        455,
        175,
        389,
        61,
        309,
        467,
        492,
        453,
        159,
        276,
        165,
        417,
        173,
        157,
        12,
        209,
        269,
        36,
        226,
        356,
        92,
        267,
        482,
        318,
        219,
        119,
        176,
        245,
        74,
        13,
        450,
        196,
        215,
        28,
        167,
        366,
        442,
        201,
        341,
        68,
        2,
        484,
        328,
        44,
        423,
        403,
        105,
        109,
        480,
        271,
        134,
        336,
        299,
        148,
        365,
        135,
        277,
        87,
        244,
        301,
        218,
        59,
        280,
        283,
        55,
        499,
        133,
        316,
        407,
        146,
        56,
        394,
        386,
        297,
        285,
        137,
        58,
        214,
        142,
        6,
        124,
        48,
        60,
        212,
        75,
        50,
        412,
        458,
        127,
        45,
        266,
        202,
        368,
        138,
        260,
        41,
        193,
        88,
        114,
        410,
        95,
        382,
        416,
        281,
        434,
        359,
        98,
        462,
        300,
        352,
        230,
        247,
        117,
        64,
        287,
        405,
        224,
        19,
        259,
        305,
        220,
        150,
        477,
        111,
        448,
        78,
        103,
        7,
        385,
        151,
        429,
        325,
        273,
        317,
        470,
        454,
        170,
        223,
        5,
        307,
        396,
        315,
        53,
        154,
        446,
        24,
        255,
        227,
        76,
        456,
        250,
        321,
        330,
        391,
        355,
        49,
        479,
        387,
        216,
        39,
        251,
        312,
        217,
        136,
        262,
        322,
        344,
        466,
        242,
        100,
        388,
        38,
        323,
        376,
        379,
        279,
        239,
        85,
        306,
        181,
        485,
        120,
        333,
        334,
        17,
        395,
        81,
        374,
        147,
        139,
        185,
        42,
        1,
        424,
        199,
        225,
        113,
        438,
        128,
        338,
        156,
        493,
        46,
        160,
        11,
        3,
        171,
        464,
        62,
        238,
        431,
        440,
        302,
        65,
        308,
        348,
        125,
        174,
        195,
        77,
        392,
        249,
        82,
        350,
        444,
        232,
        186,
        494,
        384,
        275,
        129,
        294,
        246,
        357,
        102,
        96,
        73,
        15,
        263,
        296,
        236,
        29,
        340,
        152,
        149,
        143,
        437,
        172,
        190,
        34,
        158,
        254,
        295,
        483,
        397,
        337,
        72,
        343,
        178,
        404,
        270,
        346,
        205,
        377,
        486,
        497,
        370,
        414,
        240,
        360,
        490,
        94,
        256,
        8,
        54,
        398,
        183,
        228,
        162,
        399,
        289,
        83,
        86,
        197,
        243,
        57,
        25,
        288,
        488,
        372,
        168,
        206,
        188,
        491,
        452,
        353,
        478,
        421,
        221,
        430,
        184,
        204,
        26,
        211,
        140,
        155,
        468,
        161,
        420,
        303,
        30,
        449,
        131,
        500,
        20,
        71,
        79,
        445,
        425,
        293,
        411,
        400,
        320,
        474,
        272,
        413,
        329,
        177,
        122,
        21,
        347,
        314,
        451,
        101,
        367,
        311,
        40,
        476,
        415,
        418,
        363,
        282,
        469,
        89,
        274,
        481,
        475,
        203,
        268,
        393,
        261,
        200,
        121,
        164,
        472,
        10,
        284,
        14,
        358,
        153,
        383,
        67,
        473,
        373,
        191,
        144,
        16,
        345,
        361,
        433,
        116,
        331,
        489,
        66,
        106,
        487,
        426,
        99,
        27,
        141,
        264,
        439,
        371,
        213,
        18,
        253,
        292,
        130,
        409,
        278,
        419,
        90,
        496,
        447,
        465,
        461,
        339,
        80,
        31,
        70,
        233,
        326,
        37,
        265,
        252,
        222,
        118,
        198,
        406,
        286,
        380,
        104,
        304,
        351,
        408,
        180,
        22,
        364,
        381,
        401,
        234,
        375,
        459,
        319,
        229,
        207,
        291,
        52,
        463,
        427,
        23,
        235,
        32,
        208,
        192,
        349,
        231,
        354,
        435,
        182,
        428,
        332,
        378,
        290,
        108,
        258,
        471,
        115,
        47,
        457,
        166,
        43,
        495,
        63,
        110,
        107,
        422,
        324
    ]`

	var arr []int

	if err := json.Unmarshal([]byte(s), &arr); err != nil {
		panic(err)
	}

	sort.Ints(arr)

	for i, v := range arr {
		fmt.Println(i, v, v == i+1)
		if v != i+1 {
			panic(fmt.Sprintf("expected %d, got %d", i+1, v))
		}
	}

}
