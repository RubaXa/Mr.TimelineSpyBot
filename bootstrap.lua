box.cfg{listen = 3318}

box.schema.create_space('timeline_projects')
box.schema.create_space('timeline_records')
box.schema.sequence.create('GID', {min = 1, start = 1})

box.space.timeline_projects:create_index('primary')
box.space.timeline_records:create_index('primary')

box.schema.user.create('timeline_bot', {password = 'timeline_bot'})
box.schema.user.grant('timeline_bot', 'read,write,execute', 'sequence', 'GID')
box.schema.user.grant('timeline_bot', 'read,write,execute', 'space', 'timeline_projects')
box.schema.user.grant('timeline_bot', 'read,write,execute', 'space', 'timeline_records')
