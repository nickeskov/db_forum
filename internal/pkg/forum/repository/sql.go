package repository

var sqlGetForumUserWithSince = map[bool]string{
	true: `
		SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND user_nickname < $2
		ORDER BY user_nickname DESC
		LIMIT $3`,

	false: `
		SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND user_nickname > $2
		ORDER BY user_nickname
		LIMIT $3`,
}

var sqlGetForumUser = map[bool]string{
	true: `
		SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		ORDER BY user_nickname DESC
		LIMIT $2`,

	false: `
		SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		ORDER BY user_nickname
		LIMIT $2`,
}
