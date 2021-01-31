package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter dataset file name:")
	var filename string
	for reader.Scan(){
		filename = reader.Text()
		break
	}
	input := readFile(filename)
	lines := strings.Split(input, "\n")
	unusedPizzas := getPizzas(lines)
	teams := getTeams(lines)
	ingreds := []string{}
	teamsIng := make(map[int]map[string]bool)
	for _, pizza := range unusedPizzas {
		ingreds = append(ingreds, pizza.Ingredients...)
	}
	used := make(map[int]bool)
	pizzaIndex := 0
	avoid:=false
	avoidedIndexes:=[]int{}
	avoided:=0
	max:=0
	for idx, team := range teams {
		teamsIng[idx] = make(map[string]bool)
		for i := 0; i < team.Members; i++ {
			if len(team.Pizzas)==team.Members{
				break
			}
			if pizzaIndex>len(unusedPizzas)-1{
				break
			}
			for j := pizzaIndex; j <= pizzaIndex+team.Members; j++ {
				if j > len(unusedPizzas)-1 {
					break
			}
			if used[j] {
				continue
			}
			for ingIdx, ing := range unusedPizzas[j].Ingredients {
				if teamsIng[idx][ing] {
					if len(avoidedIndexes)<=ingIdx{
						if max<10{
							team.Members+=1
							max++
						}
						continue
					}
					teams[idx].Pizzas = append(teams[idx].Pizzas, unusedPizzas[avoidedIndexes[len(unusedPizzas)]-j])
					used[unusedPizzas[avoidedIndexes[avoided]].Index] = true
					avoid = true
					continue
				}
				teamsIng[idx][ing] = true
			}
			if len(teams[idx].Pizzas) == teams[idx].Members {
				continue
			}
			if avoid {
				if max<10{
					team.Members+=1
					max++
				}
				avoid = false
				avoidedIndexes = append(avoidedIndexes, unusedPizzas[pizzaIndex].Index)
				continue
			}
			if used[unusedPizzas[pizzaIndex].Index] {
				if max<10{
					team.Members+=1
					max++
				}
				pizzaIndex += 1
				continue
			}
			teams[idx].Pizzas = append(teams[idx].Pizzas, unusedPizzas[pizzaIndex])
			used[unusedPizzas[pizzaIndex].Index] = true
			pizzaIndex += 1
			max=0
		}
			}

	}
	uniques := make(map[string]bool)
	for _, ing := range ingreds {
		if uniques[ing] {
			continue
		}
		uniques[ing] = true
	}
	res := []Team{}
	//idx := 0
	//tries:=make(map[int]int)
	for _, team := range teams {
		//for p:=idx; p<team.Members-len(team.Pizzas)+idx;p++{
		//	if used[avoidedIndexes[idx]]{
		//		continue
		//	}
		//	used[avoidedIndexes[idx]] = true
		//	teams[i].Pizzas = append(team.Pizzas, unusedPizzas[avoidedIndexes[idx]])
		//	avoidedIndexes = append(avoidedIndexes[:idx], avoidedIndexes[idx+1:]...)
		//	idx+=1
		//}
		//for team.Members>len(team.Pizzas) && idx>len(avoidedIndexes)-1{
		//	team.Pizzas = append(team.Pizzas, unusedPizzas[avoidedIndexes[idx]])
		//	idx++
		//}
		if team.Members == len(team.Pizzas) {
			res = append(res, team)
		}
	}
	s := fmt.Sprintf("%d\n", len(res))
	for _, team := range res {
		s = s + fmt.Sprintf("%d ", team.Members)
		for _, pizza := range team.Pizzas {
			s = s + fmt.Sprintf("%d ", pizza.Index)
		}
		s += "\n"

	}
	writeFile("output.txt", s)
	fmt.Println("Total Deliveries ", len(res))
	fmt.Println("Total Points ", benchmark(unusedPizzas, res))
}

type Team struct {
	Members     int
	Pizzas      []Pizza
	Ingredients map[string]bool
}

type Pizza struct {
	Index       int
	Ingredients []string
	IsUsed      bool
}

//func (p *Pizza) hasIngredient(ingredient string) bool {
//	return p.Ingredients[ingredient]
//}

func readFile(fileName string) string {
	res, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("file not found")
	}
	return string(res)
}

func writeFile(name, content string) {
	ioutil.WriteFile(name, []byte(content), 0644)
}

func getPizzas(lines []string) []Pizza {
	lineZero := strings.Split(lines[0], " ")
	noPizza, _ := strconv.Atoi(lineZero[0])
	pizzas := []Pizza{}
	for i := 1; i < noPizza; i++ {
		if len(lines) < i {
			break
		}
		pizza := Pizza{}
		ingredients := strings.Split(lines[i], " ")
		//pizza.Ingredients = make(map[string]bool)
		for j := 1; j < len(ingredients); j++ {
			pizza.Ingredients = append(pizza.Ingredients, ingredients[j])
			//pizza.Ingredients[ingredients[j]] = true
		}
		pizza.Index = i - 1
		pizza.IsUsed = false
		pizzas = append(pizzas, pizza)
	}
	return pizzas
}
func benchmark(pizzas []Pizza, solution []Team) int {
	used := make(map[int]bool)
	teamsIng := make(map[int]map[string]bool)
	res := make(map[int]int)
	for idx, team := range solution {
		teamsIng[idx] = make(map[string]bool)
		for _, pizza := range team.Pizzas {
			if used[pizza.Index] {
				continue
			}
			for _, ing := range pizza.Ingredients {
				if teamsIng[idx][ing] {
					continue
				}
				used[pizza.Index] = true
				teamsIng[idx][ing] = true
				res[idx]++
			}
		}
	}
	mark := 0
	for _, v := range res {
		mark += v * v
	}
	return mark
}

func getTeams(lines []string) []Team {
	lineZero := strings.Split(lines[0], " ")
	noT2, _ := strconv.Atoi(lineZero[1])
	noT3, _ := strconv.Atoi(lineZero[2])
	noT4, _ := strconv.Atoi(lineZero[3])
	teams := []Team{}
	teamsCount := make(map[int]int)
	teamsCount[2] = noT2
	teamsCount[3] = noT3
	teamsCount[4] = noT4
	for k, v := range teamsCount {
		for i := 1; i <= v; i++ {
			teams = append(teams, Team{Members: k})
		}
	}
	return teams
}
