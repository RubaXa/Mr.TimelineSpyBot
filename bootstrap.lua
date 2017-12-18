box.cfg{listen = 3318}

box.schema.create_space('timeline_projects')
box.schema.create_space('timeline_tokens')
box.schema.create_space('timeline_records')
box.schema.sequence.create('GID', {min = 1, start = 1})

box.space.timeline_projects:create_index('primary')

box.space.timeline_tokens:create_index('primary')
box.space.timeline_tokens:create_index('token', {type = 'hash', parts = {3, 'string'}})

box.space.timeline_records:create_index('primary')
box.space.timeline_records:create_index('project', {type = 'TREE', unique = false,  parts = {2, 'unsigned'}})

box.schema.user.create('timeline_bot', {password = 'timeline_bot'})
box.schema.user.grant('timeline_bot', 'read,write,execute', 'universe')
