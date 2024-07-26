# WebSocket API

Use [centrifuge client SDK API](https://centrifugal.dev/docs/transports/client_api) to connect to a WebSocket endpoint.


## Server Events

### `new_comment`
This event is broadcasted by the server to all clients subscribed to the channel when a new comment is created.

#### Payload
```json
{
    "payload": {
        "id": string,
        "post_id": string,
        "user": {
            "id": string,
            "first_name": string,
            "last_name": string,
            "avatar_url": string,
        },
        "body": string,
        "created_at": string,
        "updated_at": string,
    }
}
```
### `edit_comment`
This event is broadcasted by the server to all clients subscribed to the channel when a comment is updated.

#### Payload
```json
{
    "payload": {
        "id": string,
        "post_id": string,
        "user": {
            "id": string,
            "first_name": string,
            "last_name": string,
            "avatar_url": string,
        },
        "body": string,
        "created_at": string,
        "updated_at": string,
    }
}
```
### `remove_comment`
This event is broadcasted by the server to all clients subscribed to the channel when a comment is deleted.

#### Payload
```json
{
    "payload": {
        "id": string,
    }
}
```

## Client Events

### `create_comment`
This event is send by the client to create a comment. After receiving this event, the server will broadcast the comment to all clients subscribed to the channel.

#### Payload
```json
{
    "payload": {
        "post_id": string,
        "body": string,
    }
}
```

### `update_comment`
This event is send by the client to update a comment. After receiving this event, the server will broadcast the updated comment to all clients subscribed to the channel. Only the author of the comment can update it.

#### Payload
```json
{
    "payload": {
        "comment_id": string,
        "body": string,
    }
}
```

### `delete_comment`
This event is send by the client to delete a comment. After receiving this event, the server will broadcast the deleted comment to all clients subscribed to the channel. Only the author of the comment can delete it. 

#### Payload
```json
{
    "payload": {
        "comment_id": string,
    }
}
```