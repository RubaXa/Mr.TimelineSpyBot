Mr.TimelineBot
--------------


### Tarantool

```
box.schema.space.create('timeline');
box.space.timeline:create_index('pk');

box.schema.user.create('tm_bot');
box.schema.user.passwd('tm_bot', 'tm_bot')

box.schema.user.grant('tm_bot', 'read,write', 'space', 'timeline');
```