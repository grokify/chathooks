Adding Userlike Notifications to Glip
=====================================

## Events

The following is a list of events from [Userlike's API addon page](https://www.userlike.com/en/public/tutorial/addon/api).

| TYPE | EVENT | DESCRIPTION |
|------|-------|-------------|
| `offline_message` | `receive` | Receive a callback for each new offline message you receive. |
| `chat_meta` | `start` | Receive a callback for each new chat session. |
| `chat_meta` | `forward` | Receive a callback when are chat session gets forwarded. |
| `chat_meta` | `rating` | Receive a callback when a chat session receives a rating. |
| `chat_meta` | `feedback` | Receive a callback when a chat session receives a feedback. |
| `chat_meta` | `survey` | Receive a callback when a chat session receives a survey. |
| `chat_meta` | `receive` | Receive a callback when a chat session ends and the conversation is finished. |
| `chat_meta` | `goal` | Receive a callback when a goal was reached. |
| `chat_widget` | `config` | Receive a callback when a chat widget configuration changes. |
| `operator` | `online` | Receive a callback when an operator goes online. |
| `operator` | `offline` | Receive a callback when an operator goes offline. |
| `operator` | `away` | Receive a callback when an operator goes away. |
| `operator` | `back` | Receive a callback when an operator comes back. |
