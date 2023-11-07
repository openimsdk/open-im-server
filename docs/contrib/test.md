# OpenIM RPC Service Test Control Script Documentation

This document serves as a comprehensive guide to understanding and utilizing the `test.sh` script for testing OpenIM RPC services. The `test.sh` script is a collection of bash functions designed to test various aspects of the OpenIM RPC services, ensuring that each part of the API is functioning as expected.

+ Scripts：https://github.com/OpenIMSDK/Open-IM-Server/tree/main/scripts/install/test.sh

For some complex, bulky functional tests, performance tests, and various e2e tests, We are all in the current warehouse to https://github.com/OpenIMSDK/Open-IM-Server/tree/main/test or https://github.com/openim-sigs/test-infra directory In the.

+ About OpenIM Feature [Test Docs](https://docs.google.com/spreadsheets/d/1zELWkwxgOOZ7u5pmYCqqaFnvZy2SVajv/edit?usp=sharing&ouid=103266350914914783293&rtpof=true&sd=true)


## Usage

The `test.sh` script is located within the `./scripts/install/` directory of the OpenIM service's codebase. To use the script, navigate to this directory from your terminal:

```bash
cd ./scripts/install/
chmod +x test.sh
```

### Running the Entire Test Suite

To execute all available tests, you can either call the script directly or use the `make` command:

```
./test.sh openim::test::test
```

Or, if you have a `Makefile` that defines the `test-api` target:

```bash
make test-api
```

Alternatively, you can invoke specific test functions by passing them as arguments:

```
./test.sh openim::test::<function_name>
```

This `make` command should be equivalent to running `./test.sh openim::test::test`, provided that the `Makefile` is configured accordingly.



### Executing Individual Test Functions

If you wish to run a specific set of tests, you can call the relevant function by passing it as an argument to the script. Here are some examples:

**Message Tests:**

```bash
./test.sh openim::test::msg
```

**Authentication Tests:**

```bash
./test.sh openim::test::auth
```

**User Tests:**

```bash
./test.sh openim::test::user
```

**Friend Tests:**

```bash
./test.sh openim::test::friend
```

**Group Tests:**

```bash
./test.sh openim::test::group
```

Each of these commands will run the test suite associated with the specific functionality of the OpenIM service.



### Detailed Function Test Examples

T**esting Message Sending and Receiving:**

To test message functionality, the `openim::test::msg` function is called. It will register a user, send a message, and clear messages to ensure that the messaging service is operational.

```bash
./test.sh openim::test::msg
```

**Testing User Registration and Account Checks:**

The `openim::test::user` function will create new user accounts and perform a series of checks on these accounts to verify that user registration and account queries are functioning properly.

```bash
./test.sh openim::test::user
```

**Testing Friend Management:**

By invoking `openim::test::friend`, the script will test adding friends, checking friendship status, managing friend requests, and handling blacklisting.

```bash
./test.sh openim::test::friend
```

**Testing Group Operations:**

The `openim::test::group` function tests group creation, member addition, information retrieval, and member management within groups.

```bash
./test.sh openim::test::group
```

### Log Output

Each test function will output logs to the terminal to confirm the success or failure of the tests. These logs are crucial for identifying issues and verifying that each part of the service is tested thoroughly.

Each function logs its success upon completion, which aids in debugging and understanding the test flow. The success message is standardized across functions:

```
openim::log::success "<Test suite name> completed successfully."
```

By following the guidelines and instructions outlined in this document, you can effectively utilize the `test.sh` script to test and verify the OpenIM RPC services' functionality.



## Function feature

| Function Name                                        | Corresponding API/Action                      | Function Purpose                                             |
| ---------------------------------------------------- | --------------------------------------------- | ------------------------------------------------------------ |
| `openim::test::msg`                                  | Messaging Operations                          | Tests all aspects of messaging, including sending, receiving, and managing messages. |
| `openim::test::auth`                                 | Authentication Operations                     | Validates the authentication process and session management, including token handling and forced logout. |
| `openim::test::user`                                 | User Account Operations                       | Covers testing for user account creation, retrieval, updating, and overall management. |
| `openim::test::friend`                               | Friend Relationship Operations                | Ensures friend management functions correctly, including requests, listing, and blacklisting. |
| `openim::test::group`                                | Group Management Operations                   | Checks group-related functionalities like creation, invitation, information retrieval, and member management. |
| `openim::test::send_msg`                             | Send Message API                              | Simulates sending a message from one user to another or within a group. |
| `openim::test::revoke_msg`                           | Revoke Message API (TODO)                     | (Planned) Will test the revocation of a previously sent message. |
| `openim::test::user_register`                        | User Registration API                         | Registers a new user in the system to validate the registration process. |
| `openim::test::check_account`                        | Account Check API                             | Checks if an account exists for a given user ID.             |
| `openim::test::user_clear_all_msg`                   | Clear All Messages API                        | Clears all messages for a given user to validate message history management. |
| `openim::test::get_token`                            | Token Retrieval API                           | Retrieves an authentication token to validate token management. |
| `openim::test::force_logout`                         | Force Logout API                              | Forces a logout for a test user to validate session control. |
| `openim::test::check_user_account`                   | User Account Existence Check API              | Confirms the existence of a test user's account.             |
| `openim::test::get_users`                            | Get Users API                                 | Retrieves a list of users to validate user query functionality. |
| `openim::test::get_users_info`                       | Get User Information API                      | Obtains detailed information for a given user.               |
| `openim::test::get_users_online_status`              | Get User Online Status API                    | Checks the online status of a user to validate presence functionality. |
| `openim::test::update_user_info`                     | Update User Information API                   | Updates a user's information to validate account update capabilities. |
| `openim::test::get_subscribe_users_status`           | Get Subscribed Users' Status API              | Retrieves the status of users that a test user has subscribed to. |
| `openim::test::subscribe_users_status`               | Subscribe to Users' Status API                | Subscribes a test user to a set of user statuses.            |
| `openim::test::set_global_msg_recv_opt`              | Set Global Message Receiving Option API       | Sets the message receiving option for a test user.           |
| `openim::test::is_friend`                            | Check Friendship Status API                   | Verifies if two users are friends within the system.         |
| `openim::test::add_friend`                           | Send Friend Request API                       | Sends a friend request from one user to another.             |
| `openim::test::get_friend_list`                      | Get Friend List API                           | Retrieves the friend list of a test user.                    |
| `openim::test::get_friend_apply_list`                | Get Friend Application List API               | Retrieves friend applications for a test user.               |
| `openim::test::get_self_friend_apply_list`           | Get Self-Friend Application List API          | Retrieves the friend applications that the user has applied for. |
| `openim::test::add_black`                            | Add User to Blacklist API                     | Adds a user to the test user's blacklist to validate blacklist functionality. |
| `openim::test::remove_black`                         | Remove User from Blacklist API                | Removes a user from the test user's blacklist.               |
| `openim::test::get_black_list`                       | Get Blacklist API                             | Retrieves the blacklist for a test user.                     |
| `openim::test::create_group`                         | Group Creation API                            | Creates a new group with test users to validate group creation. |
| `openim::test::invite_user_to_group`                 | Invite User to Group API                      | Invites a user to join a group to test invitation functionality. |
| `openim::test::transfer_group`                       | Group Ownership Transfer API                  | Tests the transfer of group ownership from one member to another. |
| `openim::test::get_groups_info`                      | Get Group Information API                     | Retrieves information for specified groups to validate group query functionality. |
| `openim::test::kick_group`                           | Kick User from Group API                      | Simulates kicking a user from a group to test group membership management. |
| `openim::test::get_group_members_info`               | Get Group Members Information API             | Obtains detailed information for members of a specified group. |
| `openim::test::get_group_member_list`                | Get Group Member List API                     | Retrieves a list of members for a given group to ensure member listing is functional. |
| `openim::test::get_joined_group_list`                | Get Joined Group List API                     | Retrieves a list of groups that a user has joined to validate user's group memberships. |
| `openim::test::set_group_member_info`                | Set Group Member Information API              | Updates the information for a group member to test the update functionality. |
| `openim::test::mute_group`                           | Mute Group API                                | Tests the ability to mute a group, disabling message notifications for its members. |
| `openim::test::cancel_mute_group`                    | Cancel Mute Group API                         | Tests the ability to cancel the mute status of a group, re-enabling message notifications. |
| `openim::test::dismiss_group`                        | Dismiss Group API                             | Tests the ability to dismiss and delete a group from the system. |
| `openim::test::cancel_mute_group_member`             | Cancel Mute Group Member API                  | Tests the ability to cancel mute status for a specific group member. |
| `openim::test::join_group`                           | Join Group API (TODO)                         | (Planned) Will test the functionality for a user to join a specified group. |
| `openim::test::set_group_info`                       | Set Group Information API                     | Tests the ability to update the group information, such as the name or description. |
| `openim::test::quit_group`                           | Quit Group API                                | Tests the functionality for a user to leave a specified group. |
| `openim::test::get_recv_group_applicationList`       | Get Received Group Application List API       | Retrieves the list of group applications received by a user to validate application management. |
| `openim::test::group_application_response`           | Group Application Response API (TODO)         | (Planned) Will test the functionality to respond to a group join request. |
| `openim::test::get_user_req_group_applicationList`   | Get User Requested Group Application List API | Retrieves the list of group applications requested by a user to validate tracking of user's applications. |
| `openim::test::mute_group_member`                    | Mute Group Member API                         | Tests the ability to mute a specific member within a group, disabling their ability to send messages. |
| `openim::test::get_group_users_req_application_list` | Get Group Users Request Application List API  | Retrieves a list of user requests for group applications to validate group request management. |
