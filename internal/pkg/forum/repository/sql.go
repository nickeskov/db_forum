package repository

var sqlGetForumUserWithSince = map[bool]string{
	true: `
		SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND nickname < $2
		ORDER BY user_nickname DESC
		LIMIT $3`,

	false: `
		SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND nickname > $2
		ORDER BY nickname
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
		ORDER BY nickname DESC
		LIMIT $2`,

	false: `
		SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		ORDER BY nickname
		LIMIT $2`,
}
