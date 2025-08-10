package main

/*第一行代码 package main 定义了包名。你必须在源文件中非注释的第一行指明这个文件属于哪个包，
如：package main。package main表示一个可独立执行的程序，每个 Go 应用程序都包含一个名为 main 的包。*/

/*告诉 Go 编译器这个程序需要使用 fmt 包（的函数，或其他元素），fmt 包实现了格式化 IO（输入/输出）的函数*/
import (
	"fmt"
)

/*当标识符（包括常量、变量、类型、函数名、结构字段等等）以一个大写字母开头，如：Group1，
那么使用这种形式的标识符的对象就可以被外部包的代码所使用（客户端程序需要先导入这个包），
这被称为导出（像面向对象语言中的 public）；
标识符如果以小写字母开头，则对包外是不可见的，
但是他们在整个包的内部是可见并且可用的（像面向对象语言中的 protected ）。
*/
func main() {

	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 1; i <= 5; i++ {
		//TIP <p>To start your debugging session, right-click your code in the editor and select the Debug option.</p> <p>We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>.</p>
		fmt.Println("i =", 100/i)
	}
}
