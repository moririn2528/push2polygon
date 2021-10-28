package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/variety-jones/polygon"
)

const (
	MAX_FILE_SIZE = 100000
)

var (
	solution_tags = map[string]string{"AC": "OK", "TLE": "TO", "WA": "WA"}
)

func file_exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func read_file(file_name string) (string, error) {
	file, err := os.Open(file_name)
	if err != nil {
		log.Printf("Error: read_file, cannot open the file, %v, %v", file_name, err)
		return "", err
	}
	file_info, err := file.Stat()
	if err != nil {
		log.Printf("Error: read_file, cannot find file info, %v, %v", file_name, err)
		return "", err
	}
	if MAX_FILE_SIZE <= file_info.Size() {
		log.Printf("Error: read_file, file size is too large, %v", file_name)
		return "", err
	}
	buf := make([]byte, MAX_FILE_SIZE)
	n, err := file.Read(buf)
	if err != nil {
		log.Printf("Error: read_file, cannot read file, %v, %v", file_name, err)
		return "", err
	}
	err = file.Close()
	if err != nil {
		log.Printf("Error: read_file, cannot close file, %v, %v", file_name, err)
		return "", err
	}
	return string(buf[0:n]), nil
}

func write_file(file_name, context string) error {
	var file *os.File
	var err error
	file, err = os.Create(file_name)
	if err != nil {
		log.Printf("Error: write_file, cannot open the file, %v, %v", file_name, err)
		return err
	}
	_, err = file.WriteString(context)
	if err != nil {
		log.Printf("Error: write_file, cannot write on the file, %v, %v", file_name, err)
		return err
	}
	err = file.Close()
	if err != nil {
		log.Printf("Error: write_file, cannot close file, %v, %v", file_name, err)
		return err
	}
	return nil
}

func upload_file(api *polygon.PolygonApi, file_path, file_name string) error {
	context, err := read_file(file_path)
	if err != nil {
		log.Printf("Error: upload_file, read_file, %v, %v", file_path, err)
		return err
	}
	err = api.ProblemSaveFile(map[string]string{
		"type": "source",
		"name": file_name,
		"file": context,
	})
	if err != nil {
		log.Printf("Error: cannot upload file, %v, %v", file_name, err)
		return err
	}
	return nil
}

func upload(api *polygon.PolygonApi, problem_name, score_str string) error {
	fmt.Println("problem: " + problem_name)
	// upload solutions
	fmt.Print("upload solutions ... ")
	dirs, err := os.ReadDir(problem_name)
	if err != nil {
		log.Printf("Error: upload, ReadDir, %v", err)
	}
	for _, v := range dirs {
		if !v.IsDir() {
			continue
		}
		file_name := v.Name()
		tag := ""
		if file_name == "answer" {
			tag = "MA"
		}
		for key, val := range solution_tags {
			if strings.Contains(file_name, key) {
				tag = val
			}
		}
		if tag == "" {
			continue
		}
		context, err := read_file(problem_name + "/" + file_name + "/" + "main.cpp")
		if err != nil {
			log.Printf("Error: upload, read_file, %v, %v", file_name, err)
			return err
		}
		err = api.ProblemSaveSolution(map[string]string{
			"name": file_name + ".cpp",
			"file": context,
			"tag":  tag,
		})
		if err != nil {
			log.Printf("Error: cannot upload file, %v, %v", file_name, err)
			return err
		}
	}
	fmt.Println("finished")

	// upload validator
	path := problem_name + "/tests/validator.cpp"
	if file_exists(path) {
		fmt.Print("upload validator ... ")
		err = upload_file(api, path, "validator.cpp")
		if err != nil {
			log.Printf("Error: upload, upload_file, %v", err)
			return err
		}
		err = api.ProblemSetValidator(map[string]string{
			"validator": "validator.cpp",
		})
		if err != nil {
			log.Printf("Error: cannot set validator, %v", err)
			return err
		}
		fmt.Println("finished")
	} else {
		log.Println("Warning: No validator")
	}

	// upload checker
	path = problem_name + "/tests/output_checker.cpp"
	if file_exists(path) {
		fmt.Print("upload checker ... ")
		err = upload_file(api, path, "checker.cpp")
		if err != nil {
			log.Printf("Error: upload, upload_file, %v", err)
			return err
		}
		name := "checker.cpp"
		err = api.ProblemSetChecker(map[string]string{
			"checker": name,
		})
		if err != nil {
			log.Printf("Error: cannot set validator, %v", err)
			return err
		}
		fmt.Println("finished")
	} else {
		log.Println("Warning: No checker")
	}

	// statements cannot push japanese text
	// // upload statements
	// path = problem_name + "/statement.tex"
	// if file_exists(path) {
	// 	fmt.Print("upload statements ... ")
	// 	context, err := read_file(path)
	// 	if err != nil {
	// 		log.Printf("Error: upload, read_file, %v, %v", path, err)
	// 		return err
	// 	}
	// 	context_lines := strings.Split(context, "\n")
	// 	var score_str, title, statement, input_format, output_format, note string
	// 	str := &title
	// 	for _, v := range context_lines {
	// 		switch v {
	// 		case "<score>":
	// 			str = &score_str
	// 		case "<title>":
	// 			str = &title
	// 		case "<problem>":
	// 			str = &statement
	// 		case "<input>":
	// 			str = &input_format
	// 		case "<output>":
	// 			str = &output_format
	// 		case "<note>":
	// 			str = &note
	// 		default:
	// 			log.Println(str)
	// 			log.Println(v)
	// 			*str += v
	// 		}
	// 	}
	// 	statement_map := map[string]string{
	// 		"lang":     "Japanese",
	// 		"encoding": "UTF-8",
	// 		"name":     strings.ReplaceAll(title, "\n", ""),
	// 		"legend":   statement,
	// 		"input":    input_format,
	// 		"output":   output_format,
	// 		"notes":    note,
	// 	}
	// 	score_str = strings.TrimSpace(score_str)
	// 	if score_str != "" {
	// 		_, err := strconv.Atoi(score_str)
	// 		if err != nil {
	// 			log.Printf("Error: upload, statement.tex, score err, %v", err)
	// 			return err
	// 		}
	// 		statement_map["scoring"] = score_str
	// 	}
	// 	err = api.ProblemSaveStatement(statement_map)
	// 	if err != nil {
	// 		log.Printf("Error: cannot set statements, %v", err)
	// 		return err
	// 	}
	// 	fmt.Println("finished")
	// } else {
	// 	log.Println("Warning: No statements")
	// }

	// upload tests
	fmt.Printf("create tests ... ")
	test_count := 0
	err = api.ProblemEnableGroups(map[string]string{
		"testset": "tests",
		"enable":  "true",
	})
	if err != nil {
		log.Printf("Error: upload, ProblemEnableGroups, %v", err)
		return err
	}
	err = api.ProblemEnablePoints(map[string]string{
		"enable": "true",
	})
	if err != nil {
		log.Printf("Error: upload, ProblemEnableGroups, %v", err)
		return err
	}

	dirs, err = os.ReadDir(problem_name + "/tests")
	if err != nil {
		log.Printf("Error: upload, ReadDir, %v", err)
	}
	scored := false
	for _, v := range dirs {
		name := v.Name()
		if !strings.Contains(name, ".in") {
			continue
		}
		context, err := read_file(problem_name + "/tests/" + name)
		if err != nil {
			log.Printf("Error: upload, read_file, %v, %v, %v", problem_name, name, err)
			return err
		}
		test_count++
		m := map[string]string{
			"testset":   "tests",
			"testIndex": strconv.Itoa(test_count),
			"testInput": context,
			"testGroup": "A",
		}
		if !scored {
			scored = true
			m["testPoints"] = score_str
		}
		if strings.HasPrefix(name, "sample") {
			m["testUseInStatements"] = "true"
		}
		err = api.ProblemSaveTest(m)
		if err != nil {
			log.Printf("Error: upload, ProblemSaveTest, %v, %v", name, err)
			return err
		}
	}

	err = api.ProblemSaveTestGroups(map[string]string{
		"testset":        "tests",
		"group":          "A",
		"pointsPolicy":   "COMPLETE_GROUP",
		"feedbackPolicy": "COMPLETE",
	})
	if err != nil {
		log.Printf("Error: upload, ProblemSaveTestGroups, %v", err)
		return err
	}
	fmt.Println("finished")

	// upload generator
	path = problem_name + "/tests/generator.cpp"
	indices := []string{}
	if file_exists(path) {
		fmt.Print("upload generator ... ")
		context, err := read_file(path)
		if err != nil {
			log.Printf("Error: upload, read_file, %v", err)
			return err
		}
		context_lines := strings.Split(context, "\n")
		paren_cnt := 0
		context_blocks := make([][]string, 0)
		push_names := make([]string, 0)
		main_index := -1
		block := make([]string, 0)

		for i := 0; i < len(context_lines); i++ {
			v := context_lines[i]
			var l, r int
			cnt := paren_cnt
			paren_cnt += strings.Count(v, "{")
			paren_cnt -= strings.Count(v, "}")
			if cnt != 0 || strings.Contains(v, ";") || len(v) <= 2 {
				block = append(block, v)
				continue
			}
			if strings.Contains(v, "(") && strings.Contains(v, ")") && (strings.Contains(v, "{") ||
				(i+1 < len(context_lines) && context_lines[i+1] == "{")) {
				l = i
			} else {
				block = append(block, v)
				continue
			}
			cnt = 0
			for r = l; r < len(context_lines); r++ {
				cnt += strings.Count(context_lines[r], "{")
				cnt -= strings.Count(context_lines[r], "}")
				if cnt == 0 {
					break
				}
			}
			if len(context_lines) <= r {
				log.Println("Error: upload, generator, len(context_lines) <= r")
				return errors.New("index error")
			}

			if len(block) > 0 {
				context_blocks = append(context_blocks, block)
				push_names = append(push_names, "")
			}
			block = make([]string, 0)
			func_name := strings.TrimSpace(v[:strings.Index(v, "(")])
			func_name_l := strings.LastIndex(func_name, " ")
			func_name = func_name[func_name_l+1:]
			if func_name == "main" {
				main_index = i
				break
			}

			ofs_name := ""
			for j := l; j <= r; j++ {
				u := context_lines[j]
				if strings.Contains(u, "ofstream") {
					ofs_name_l := strings.Index(u, "ofstream") + 8
					ofs_name_r := len(u) - 1
					ofs_name = u[ofs_name_l:ofs_name_r]
					if strings.Contains(ofs_name, "(") {
						ofs_name_r = strings.Index(ofs_name, "(")
						ofs_name = ofs_name[:ofs_name_r]
					}
					continue
				}
				if ofs_name != "" {
					if strings.Contains(context_lines[j], ofs_name+".") {
						continue
					}
					u = strings.Replace(u, ofs_name, "cout", -1)
				}
				block = append(block, u)
			}

			context_blocks = append(context_blocks, block)
			push_names = append(push_names, func_name)
			i = r
			paren_cnt = 0
			block = make([]string, 0)
		}
		if main_index == -1 {
			log.Printf("Error: upload, generator, cannot find main function, %v", problem_name)
			return errors.New("main index")
		}

		gen_func_str, err := read_file(problem_name + "/tests/gen_function.txt")
		if err != nil {
			log.Println("Error: upload, read_file, gen_function.txt")
			return err
		}
		gen_func_lines := strings.Split(gen_func_str, "\n")
		gen_func := map[string]int{}
		for _, str := range gen_func_lines {
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}
			v := strings.Split(str, ":")
			if len(v) != 2 {
				log.Println("Error: upload, gen_function.txt, syntax error")
				return errors.New("syntax error")
			}
			num, err := strconv.Atoi(strings.TrimSpace(v[1]))
			func_name := strings.TrimSpace(v[0])
			if err != nil {
				log.Printf("Error: upload, gen_function.txt, %v", err)
				return err
			}
			gen_func[func_name] = num
		}

		gen_script := []string{}
		for i, b := range context_blocks {
			v := b[0]
			if push_names[i] == "" {
				continue
			}
			num, ok := gen_func[push_names[i]]
			if !ok || num == 0 {
				continue
			}
			push_context := []string{}
			for j := 0; j < i; j++ {
				push_context = append(push_context, strings.Join(context_blocks[j], "\n"))
			}
			arg_l := strings.Index(v, "(")
			arg_r := strings.LastIndex(v, ")")
			if arg_r <= arg_l+1 {
				log.Printf("Error: upload, Index, generator, arg_r <= arg_l+1, %v, %v", i, v)
				return errors.New("argment error")
			}
			args := strings.Split(v[arg_l+1:arg_r], ",")
			push_context = append(push_context, v[:arg_l+1]+v[arg_r:])
			l := 1
			if i+1 < len(context_lines) && context_lines[i+1] == "{" {
				push_context = append(push_context, "{")
				l++
			}
			for _, ar := range args {
				ar = strings.TrimSpace(strings.Replace(ar, "const", "", -1))
				push_context = append(push_context, "\t"+ar+";")
			}
			for j := l; j < len(b); j++ {
				push_context = append(push_context, b[j])
			}
			for j := i + 1; j < len(context_blocks); j++ {
				push_context = append(push_context, strings.Join(context_blocks[j], "\n"))
			}
			push_context = append(push_context, "int main(int argc, char* argv[]){\n"+
				"\tregisterGen(argc, argv, 1);\n"+
				"\t"+push_names[i]+"();\n"+
				"}\n")

			err = api.ProblemSaveFile(map[string]string{
				"type": "source",
				"name": push_names[i] + ".cpp",
				"file": strings.Join(push_context, "\n"),
			})
			if err != nil {
				log.Printf("Error: cannot upload file, generator, %v, %v", push_names[i], err)
				return err
			}
			gen_script = append(gen_script, "<#list 1.."+strconv.Itoa(num)+" as i>",
				push_names[i]+" ${i} 475 > $",
				"</#list>")
			for i := 0; i < num; i++ {
				indices = append(indices, strconv.Itoa(test_count+i+1))
			}
			test_count += num
		}
		err = api.ProblemSaveScript(map[string]string{
			"testset": "tests",
			"source":  strings.Join(gen_script, "\n"),
		})
		if err != nil {
			log.Printf("Error: upload, write_file, gen_script.txt, %v", err)
			return err
		}
		fmt.Println("finished")
	} else {
		log.Println("Warning: No generator")
	}

	err = api.ProblemSetTestGroup(map[string]string{
		"testset":     "tests",
		"testGroup":   "A",
		"testIndices": strings.Join(indices, ","),
	})
	if err != nil {
		log.Printf("Error: upload, ProblemSetTestGroup, %v", err)
		return err
	}

	return nil
}

func upload_by_name(name string) error {
	if !file_exists(name + "/polygon.txt") {
		log.Printf("Error: upload_by_name, file not exists, %v", name)
		return errors.New("file not exists")
	}
	context, err := read_file(name + "/polygon.txt")
	if err != nil {
		log.Println("Error: upload_by_name, read_file")
		return err
	}
	context_lines := strings.Split(context, "\n")
	api := polygon.PolygonApi{
		ApiKey:    "4e5716bdaf99f3144ab82557d093cc7ba4a6ef4f",
		Secret:    "3eeaa6645009b21a5b8388e56189ead0ec93e5b2",
		ProblemId: "",
	}
	score := ""
	for _, v := range context_lines {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		vs := strings.Split(v, ":")
		if len(vs) != 2 {
			continue
		}
		sets := strings.TrimSpace(vs[0])
		num_str := strings.TrimSpace(vs[1])
		_, err = strconv.Atoi(num_str)
		if err != nil {
			log.Printf("Error: upload_by_name, polygon.txt, syntax error, %v", err)
			return err
		}
		if sets == "score" {
			score = num_str
		}
		if sets == "problem_id" {
			api.ProblemId = num_str
		}
	}
	if api.ProblemId == "" || score == "" {
		log.Println("Error: upload_by_name, polygon.txt, syntax error")
		return errors.New("syntax error")
	}
	return upload(&api, name, score)
}

func upload_all() error {
	dirs, err := os.ReadDir(".")
	if err != nil {
		log.Println("Error: main, ReadDir")
		return err
	}
	for _, d := range dirs {
		name := d.Name()
		if !file_exists(name+"/polygon.txt") || name == "template" {
			continue
		}
		err := upload_by_name(name)
		if err != nil {
			log.Printf("Error: upload_all, %v", err)
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		upload_all()
		return
	}
	for _, name := range args {
		err := upload_by_name(name)
		if err != nil {
			return
		}
	}
}
