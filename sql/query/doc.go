/*
Package query is a helper package for building sql statement with fluent methods

	func testQuery() {
		sqlStr, params := query.Builder().
			Table("users as u").
			Select("username", "age", "id", "ent.name as enterprise_name", query.Raw("count(role_id) as role_count"), query.Raw("sum(login_times) * ?", 1.5)).
			WhereGroup(func(builder query.Condition) {
				builder.Where("u.age", ">", 40).
					Where("u.username", "=", "zhangsanfeng").
					Where("ent.name", "like", "重庆%")
			}).
			OrWhereGroup(func(builder query.Condition) {
				builder.Where("age", ">", 14).
					Where("username", "=", "lixiaoyao").
					WhereGroup(func(builder query.Condition) {
						builder.Where("sex", "=", "boy")
					})
			}).
			WhereIn("role_id", query.Builder().Table("roles as rol").Select("id").Where("status", "=", 1)).
			WhereIn("role_id", 13, 41, 56, 52).
			OrWhereNotIn("org_id", 14, 52).
			WhereExist(query.Builder().Table("enterprise").Select("id").WhereIn("status", 1, 2, 3)).
			WhereRaw("password=?", "asfhjq9u").
			GroupBy("username", "org_id").
			Having(func(builder query.Condition) {
				builder.WhereRaw("COUNT(id) > 10")
			}).
			LeftJoin("enterprise as ent", func(builder query.Condition) {
				builder.WhereColumn("ent.id", "=", "u.enterprise_id")
			}).
			OrderBy("created_at", "DESC").
			OrderBy("role_id", "ASC").
			Limit(1000).
			Offset(50).
			Union(query.Builder().
				Table("users_2 AS u").
				Select("username", "age", "id", "ent.name as enterprise_name", query.Raw("count(role_id) as role_count"), query.Raw("sum(login_times) * ?", 1.5)).
				LeftJoin("enterprise AS ent", func(builder query.Condition) {
					builder.WhereColumn("enterprise_id", "=", "ent.id")
				}), true).
			ResolveQuery()

		fmt.Println(replaceQuestionMark(sqlStr, params), ";")
		fmt.Println()

		sqlStr, params = query.Builder().Table("user").ResolveInsert(query.KV{
			"username":   "guan",
			"age":        24,
			"created_at": time.Now(),
		})
		fmt.Println(replaceQuestionMark(sqlStr, params), ";")
		fmt.Println()

		sqlStr, params = query.Builder().Table("user").Where("id", "=", 100).When(func() bool {
			return true
		}, func(builder query.Condition) {
			builder.WhereNull("deleted_at")
		}).OrWhen(func() bool {
			return false
		}, func(builder query.Condition) {
			builder.WhereNotNull("deleted_at")
		}).ResolveUpdate(query.KV{
			"username":   "guan",
			"age":        24,
			"created_at": time.Now(),
			"count":      query.Raw("count + 1"),
		})
		fmt.Println(replaceQuestionMark(sqlStr, params), ";")
		fmt.Println()

		sqlStr, params = query.Builder().Table("user").Where("id", "=", 100).OrderBy("created_at", "DESC").Limit(5).ResolveDelete()
		fmt.Println(replaceQuestionMark(sqlStr, params), ";")
		fmt.Println()
	}

	func replaceQuestionMark(sql string, values []interface{}) string {
		fmt.Println(sql)
		fmt.Println(values)
		return fmt.Sprintf(strings.ReplaceAll(sql, "?", "'%v'"), values...)
	}
*/
package query