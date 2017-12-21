--cfg--
box.cfg{listen = 3318}
box.schema.sequence.create('GID', {min = 1, start = 1})

--Projects--
box.schema.create_space('timeline_projects')
box.space.timeline_projects:create_index('primary')

--Tokens--
box.schema.create_space('timeline_tokens')
box.space.timeline_tokens:create_index('primary')
box.space.timeline_tokens:create_index('token', {type = 'hash', parts = {3, 'string'}})

--Records--
box.schema.create_space('timeline_records')
box.space.timeline_records:create_index('primary')
box.space.timeline_records:create_index('project', {type = 'TREE', unique = false,  parts = {2, 'unsigned'}})

--User--
box.schema.user.create('timeline_bot', {password = 'timeline_bot'})
box.schema.user.grant('timeline_bot', 'read,write,execute', 'universe')
