package main

func main() {
	program := `
fun MoreThan50(x:int) -> bool {
let x:int = 23; //syntax ok, but this should not be allowed!!
if (x <= 50) {
return false;
}
return true;
}

let x:int = 45; //this is fine
while (x < 50) {
__print MoreThan50(x); //"false" *5 since bool operator is <
x = x + 1;
}

let x:int = 45; //re-declaration in the same scope ... not allowed!!
while (MoreThan50(x)) {
__print MoreThan50(x); //"false" x5 since bool operator is <=
x = x + 1;
}

let w: int = __width;
let h: int = __height;

for (let u:int = 0; u<w; u = u+1)
{
for (let v:int = 0; v<h; v = v+1)
{
//set the pixel at u,v to the colour green
__write_box u,v,1,1,#00ff00;
//or ... assume one pixel 1x1
//__write u,v,#00ff00;
{}{}{}
}
}
	`

	parser := NewParser(program)
	printVisitor := PrintNodesVisitor{}
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		panic(err)
	}
	node.Accept(&printVisitor)
}
