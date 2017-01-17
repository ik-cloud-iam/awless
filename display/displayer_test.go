package display

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
	"github.com/wallix/awless/rdf"
)

func TestTabularDisplays(t *testing.T) {
	g := createIntancesGraph()
	headers := []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
	}
	displayer := BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance,
		Format:  "csv",
	})
	displayer.SetGraph(g)

	expected := `Id, Name, State, Type, Public IP
inst_1, redis, running, t2.micro, 1.2.3.4
inst_2, django, stopped, t2.medium, 
inst_3, apache, running, t2.xlarge, `
	var w bytes.Buffer
	err := displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance,
		Format:  "csv",
		SortBy:  []string{"Name"},
	})
	displayer.SetGraph(g)

	expected = `Id, Name, State, Type, Public IP
inst_3, apache, running, t2.xlarge, 
inst_2, django, stopped, t2.medium, 
inst_1, redis, running, t2.micro, 1.2.3.4`

	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	headers = []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "State"},
			ColoredValues:          map[string]color.Attribute{"running": color.FgGreen},
		},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
	}
	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
	})
	displayer.SetGraph(g)
	expected = `+--------+--------+---------+-----------+-----------+
|  ID ▲  |  NAME  |  STATE  |   TYPE    | PUBLIC IP |
+--------+--------+---------+-----------+-----------+
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_2 | django | stopped | t2.medium |           |
| inst_3 | apache | running | t2.xlarge |           |
+--------+--------+---------+-----------+-----------+
`
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy: []string{"state", "id"},
	})
	displayer.SetGraph(g)
	expected = `+--------+--------+---------+-----------+-----------+
|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
+--------+--------+---------+-----------+-----------+
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_3 | apache | running | t2.xlarge |           |
| inst_2 | django | stopped | t2.medium |           |
+--------+--------+---------+-----------+-----------+
`
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy: []string{"state", "name"},
	})
	displayer.SetGraph(g)
	expected = `+--------+--------+---------+-----------+-----------+
|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
+--------+--------+---------+-----------+-----------+
| inst_3 | apache | running | t2.xlarge |           |
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_2 | django | stopped | t2.medium |           |
+--------+--------+---------+-----------+-----------+
`
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	headers = []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
	}
	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "porcelain",
	})
	displayer.SetGraph(g)
	expected = `inst_1
redis
inst_2
django
inst_3
apache`
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestDateLists(t *testing.T) {
	users := []byte(`/region<eu-west-1>	"has_type"@[]	"/region"^^type:text
/region<eu-west-1>	"parent_of"@[]	/user<user1>
/region<eu-west-1>	"parent_of"@[]	/user<user2>
/region<eu-west-1>	"parent_of"@[]	/user<user3>
/user<user1>	"has_type"@[]	"/user"^^type:text
/user<user2>	"has_type"@[]	"/user"^^type:text
/user<user3>	"has_type"@[]	"/user"^^type:text
/user<user1>	"property"@[]	"{"Key":"Id","Value":"user1"}"^^type:text
/user<user2>	"property"@[]	"{"Key":"Id","Value":"user2"}"^^type:text
/user<user3>	"property"@[]	"{"Key":"Id","Value":"user3"}"^^type:text
/user<user1>	"property"@[]	"{"Key":"Name","Value":"my_username_1"}"^^type:text
/user<user2>	"property"@[]	"{"Key":"Name","Value":"my_username_2"}"^^type:text
/user<user3>	"property"@[]	"{"Key":"Name","Value":"my_username_3"}"^^type:text
/user<user2>	"property"@[]	"{"Key":"PasswordLastUsedDate","Value":"2016-12-22T11:13:23Z"}"^^type:text
/user<user3>	"property"@[]	"{"Key":"PasswordLastUsedDate","Value":"2016-12-10T08:35:37Z"}"^^type:text`)

	g := rdf.NewGraph()
	g.Unmarshal(users)

	headers := []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "PasswordLastUsedDate"}, Format: Short},
	}

	displayer := BuildGraphDisplayer(headers, Options{
		RdfType: rdf.User, Format: "table",
		SortBy: []string{"id"},
	})
	displayer.SetGraph(g)
	expected := `+-------+---------------+----------------------+
| ID ▲  |     NAME      | PASSWORDLASTUSEDDATE |
+-------+---------------+----------------------+
| user1 | my_username_1 |                      |
| user2 | my_username_2 | 12/22/16 11:13       |
| user3 | my_username_3 | 12/10/16 08:35       |
+-------+---------------+----------------------+
`
	var w bytes.Buffer
	err := displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.User, Format: "table",
		SortBy: []string{"passwordlastuseddate"},
	})
	displayer.SetGraph(g)
	expected = `+-------+---------------+------------------------+
|  ID   |     NAME      | PASSWORDLASTUSEDDATE ▲ |
+-------+---------------+------------------------+
| user1 | my_username_1 |                        |
| user3 | my_username_3 | 12/10/16 08:35         |
| user2 | my_username_2 | 12/22/16 11:13         |
+-------+---------------+------------------------+
`
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestMaxWidth(t *testing.T) {
	g := createIntancesGraph()
	headers := []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
	}
	displayer := BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy: []string{"state", "name"},
	})
	displayer.SetGraph(g)
	expected := `+--------+--------+---------+-----------+-----------+
|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
+--------+--------+---------+-----------+-----------+
| inst_3 | apache | running | t2.xlarge |           |
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_2 | django | stopped | t2.medium |           |
+--------+--------+---------+-----------+-----------+
`
	var w bytes.Buffer
	err := displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	headers = []ColumnDefinition{
		StringColumnDefinition{Prop: "Id", TruncateSize: 4, TruncateRight: true},
		StringColumnDefinition{Prop: "Name", DontTruncate: true},
		StringColumnDefinition{Prop: "State", DontTruncate: true},
		StringColumnDefinition{Prop: "Type", TruncateSize: 6},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP", DontTruncate: true},
	}
	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy: []string{"state", "name"},
	})
	displayer.SetGraph(g)
	expected = `+------+--------+---------+--------+-----------+
|  ID  |  NAME  | STATE ▲ |  TYPE  | PUBLIC IP |
+------+--------+---------+--------+-----------+
| i... | apache | running | ...rge |           |
| i... | redis  | running | ...cro | 1.2.3.4   |
| i... | django | stopped | ...ium |           |
+------+--------+---------+--------+-----------+
`
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	headers = []ColumnDefinition{
		StringColumnDefinition{Prop: "Id", Friendly: "I", TruncateSize: 5},
		StringColumnDefinition{Prop: "Name", Friendly: "N", TruncateSize: 5},
		StringColumnDefinition{Prop: "State", Friendly: "S", TruncateSize: 5},
		StringColumnDefinition{Prop: "Type", Friendly: "T", TruncateSize: 5},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "P", TruncateSize: 5},
	}
	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy:   []string{"s", "n"},
		MaxWidth: 0,
	})
	displayer.SetGraph(g)
	expected = `+-------+-------+-------+-------+-------+
|   I   |   N   |  S ▲  |   T   |   P   |
+-------+-------+-------+-------+-------+
| ..._3 | ...he | ...ng | ...ge |       |
| ..._1 | redis | ...ng | ...ro | ....4 |
| ..._2 | ...go | ...ed | ...um |       |
+-------+-------+-------+-------+-------+
`
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy:   []string{"s", "n"},
		MaxWidth: 50,
	})
	displayer.SetGraph(g)
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy:   []string{"s", "n"},
		MaxWidth: 21,
	})
	displayer.SetGraph(g)
	expected = `+-------+-------+-------+
|   I   |   N   |  S ▲  |
+-------+-------+-------+
| ..._3 | ...he | ...ng |
| ..._1 | redis | ...ng |
| ..._2 | ...go | ...ed |
+-------+-------+-------+
Columns truncated to fit terminal: 'T', 'P'
`
	w.Reset()
	err = displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestCompareInterface(t *testing.T) {
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(4)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(1)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(-3)), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}("abc"), interface{}("bbc")), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}("abc"), interface{}("aac")), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1.2), interface{}(1.3)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1.2), interface{}(1.1)), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
}

func createIntancesGraph() *rdf.Graph {
	instances := []byte(`/region<eu-west-1> "has_type"@[] "/region"^^type:text
  /instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Type","Value":"t2.micro"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"PublicIp","Value":"1.2.3.4"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"State","Value":"running"}"^^type:text

  /instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Name","Value":"django"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Type","Value":"t2.medium"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"State","Value":"stopped"}"^^type:text
	
	
  /instance<inst_3>  "has_type"@[] "/instance"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Id","Value":"inst_3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Name","Value":"apache"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Type","Value":"t2.xlarge"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"State","Value":"running"}"^^type:text

  /region<eu-west-1>  "parent_of"@[] /instance<inst_1>
  /region<eu-west-1>  "parent_of"@[] /instance<inst_2>`)

	g := rdf.NewGraph()
	g.Unmarshal(instances)
	return g
}