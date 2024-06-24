# tractive

[Tractive](https://tractive.com/) API wrapper written in Go. Based on https://github.com/FAXES/tractive .


## Feature matrix

| Authentication                       |    |
|--------------------------------------|----|
| Authentication                       | ✅ |
| Re-authentication when token expires | ❌ |

| Account                   |    |
|---------------------------|----|
| Get account info          | ✅ |
| Get account subscriptions | ✅ |
| Get account subscription  | ✅|
| Get account shares        | ✅ |

| Commands              |    |
|-----------------------|----|
| Enable live tracking  | ❌ |
| Disable live tracking | ❌ |
| Turn LED on           | ❌ |
| Turn LED off          | ❌ |
| Turn buzzer on        | ❌ |
| Turn buzzer off       | ❌ |

| Pet      |    |
|----------|----|
| Get pet  | ✅ |
| Get pets | ✅ |

| Tracker              |    |
|----------------------|----|
| Get all trackers     | ❌ |
| Get tracker          | ❌ |
| Get tracker history  | ❌ |
| Get tracker location | ❌ |
| Get tracker hardware | ❌ |
