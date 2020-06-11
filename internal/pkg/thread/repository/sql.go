package repository

var sqlGetThreadsByForumSlugSince = map[bool]string{
	true: `
		SELECT id,
			   slug,
			   forum_slug,
			   author_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1 AND created <= $2
		ORDER BY created DESC
		LIMIT $3`,

	false: `
		SELECT id,
			   slug,
			   forum_slug,
			   author_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1 AND created >= $2
		ORDER BY created
		LIMIT $3`,
}

var sqlGetThreadsByForumSlug = map[bool]string{
	true: `
		SELECT id,
			   slug,
			   forum_slug,
			   author_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1
		ORDER BY created DESC
		LIMIT $2`,

	false: `
		SELECT id,
			   slug,
			   forum_slug,
			   author_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1
		ORDER BY created
		LIMIT $2`,
}
