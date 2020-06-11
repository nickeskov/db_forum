package repository

var sqlGetForumUser = map[bool]string{
	true: ` 
		SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND nickname < $2
		ORDER BY nickname DESC
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
