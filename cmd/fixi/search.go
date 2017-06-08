package main

// func SearchCmd(ctx *cli.Context) error {
// 	queryStr := strings.Join(ctx.Args(), " ")
// 	if queryStr == "" {
// 		return cli.ShowCommandHelp(ctx, "search")
// 	}
//
// 	fixi, err := loadFixity(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	q := q.FromString(queryStr)
//
// 	results, err := fixi.Search(q)
// 	if err != nil {
// 		return err
// 	}
//
// 	listable, err := listableResults(fixi, results)
// 	if err != nil {
// 		return err
// 	}
//
// 	w := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
// 	fmt.Fprintln(w, "\t"+strings.Join(listable.UpperFields(), "\t"))
//
// 	for i, row := range listable.Rows() {
// 		fmt.Fprintln(w, fmt.Sprintf("%d\t%s", i+1, strings.Join(row, "\t")))
// 	}
//
// 	return w.Flush()
// }
//
// // listable was just hacked together, it should be improved
// type listable struct {
// 	fields []string
// 	maps   []map[string]interface{}
// }
//
// func listableResults(fixi fixity.Fixity, hashes []string) (*listable, error) {
// 	listable := newListable()
//
// 	// always have a base set of values, regardless of the number of results
// 	listable.AddField("id")
//
// 	for _, h := range hashes {
// 		row := map[string]interface{}{}
// 		v, err := fixi.ReadHash(h)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		row["id"] = v.Id
//
// 		for _, jsonWithMeta := range v.MultiJson {
// 			var j map[string]interface{}
// 			if err := json.Unmarshal(jsonWithMeta.JsonBytes, &j); err != nil {
// 				return nil, err
// 			}
//
// 			// TODO(leeola)" figure out expected behavior for duplicate fields?
// 			for k, v := range j {
// 				listable.AddField(k)
// 				row[k] = v
// 			}
// 		}
//
// 		listable.AddRow(row)
// 	}
//
// 	return listable, nil
// }
//
// func newListable() *listable {
// 	return &listable{}
// }
//
// func (l *listable) AddField(f string) {
// 	var exists bool
// 	for _, existingF := range l.fields {
// 		if existingF == f {
// 			exists = true
// 		}
// 	}
//
// 	if !exists {
// 		l.fields = append(l.fields, f)
// 	}
// }
//
// func (l *listable) AddRow(m map[string]interface{}) {
// 	l.maps = append(l.maps, m)
// }
//
// func (l *listable) Rows() [][]string {
// 	var rows [][]string
// 	for _, m := range l.maps {
// 		var row []string
// 		for _, f := range l.fields {
// 			v, _ := m[f]
// 			var s string
// 			if v != nil {
// 				s = fmt.Sprintf("%s", v)
// 			}
// 			row = append(row, s)
// 		}
// 		rows = append(rows, row)
// 	}
// 	return rows
// }
//
// func (l *listable) UpperFields() []string {
// 	fs := make([]string, len(l.fields))
// 	for i, f := range l.fields {
// 		fs[i] = strings.ToUpper(f)
// 	}
// 	return fs
// }
