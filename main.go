package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func main() {

	gameMenu()

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func writePlayerMapToCSV(playerMap map[int]Player) {
	file, err := os.Create("players.csv")
	checkError("Error:", err)
	defer file.Close()
	headers := []string{
		"Name",
		"LifeTotal",
		"Power",
		"Mana",
		"Experience",
	}
	// write column headers
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for key := range playerMap {
		r := make([]string, 0, 1+len(headers))
		lifeString := strconv.Itoa(playerMap[key].LifeTotal)
		powerString := strconv.Itoa(playerMap[key].Power)
		manaString := strconv.Itoa(playerMap[key].Mana)
		experienceString := strconv.Itoa(playerMap[key].Experience)
		r = append(
			r,
			playerMap[key].Name,
			lifeString,
			powerString,
			manaString,
			experienceString,
		)
		writer.Write(r)
	}
}

func writePlayerMapToJSON(playerMap map[int]Player) {
	for key := range playerMap {

	}
	data := User{
		Name:     "bob",
		Password: "bob",
		Character: Player{
			Name:       "bob",
			LifeTotal:  100,
			Power:      15,
			Mana:       0,
			Experience: 45,
		},
	}
	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile("players.json", file, 0644)
}

func getPlayersFromCSV() map[int]Player {
	var cpt int
	playerMap := make(map[int]Player)
	csvfile, err := os.Open("players.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		LifeTotal, err := strconv.Atoi(record[1])
		Power, err := strconv.Atoi(record[2])
		Mana, err := strconv.Atoi(record[3])
		Experience, err := strconv.Atoi(record[4])
		player := Player{record[0], LifeTotal, Power, Mana, Experience}
		playerMap[cpt] = player
		cpt += 1
		//fmt.Printf("Name: %s LifeTotal: %s Power: %s Mana: %s Experience: %s\n", record[0], record[1], record[2], record[3], record[4])
	}
	return playerMap
}

type Player struct {
	Name       string
	LifeTotal  int
	Power      int
	Mana       int
	Experience int
}

type User struct {
	Name         string
	Password     string
	HashPassword string
	Character    Player
}

func (p *Player) hit(playerHit *Player, Power int) {
	fmt.Println(p.Name, "tries to hit ", playerHit.Name, " with a Power of ", Power)
	playerHit.getHit(Power)
	p.Mana += 1
	//fmt.Println(p.Name, " has ", p.Mana, " Mana points ")
}

func (p *Player) getHit(Power int) {
	p.LifeTotal = p.LifeTotal - Power
	//fmt.Println(p.Name, " has took a torgnole and now has ", p.LifeTotal, " HP ")
}

func (p *Player) getXp() {
	p.Experience += 1
	//fmt.Println(p.Name, " has took a torgnole and now has ", p.LifeTotal, " HP ")
}

func combat(player1 *Player, player2 *Player, key1 int, key2 int) int {
	var f int
	for f != 4 {
		if player1.LifeTotal <= 0 {
			fmt.Println(player2)
			fmt.Println(player2.Name, " has won !")
			player1.LifeTotal = 100
			player2.LifeTotal = 100
			player1.Mana = 0
			player2.Mana = 0
			player2.Experience = player2.Experience + 1
			f = 4
			return key2
		}
		if player2.LifeTotal <= 0 {
			fmt.Println(player1)
			fmt.Println(player1.Name, " has won !")
			player1.LifeTotal = 100
			player2.LifeTotal = 100
			player1.Mana = 0
			player2.Mana = 0
			player1.Experience = player1.Experience + 1
			f = 4
			return key1
		}
		fmt.Println("You're in a fight, what do you do ?")
		fmt.Println("Your life: ", player1.LifeTotal, ", your Mana : ", player1.Mana)
		fmt.Println("Your ennemy : ", player2.Name, " ", player2.LifeTotal, " HP , ", player2.Mana, " Mana points ")
		fmt.Println("1 - Simple Attack\n2 - Spell\n3 - run away ??")
		fmt.Scanf("%d", &f)
		switch f {
		case 1:
			player1.hit(player2, player1.Power)
			ennemyFightBack(player1, player2)
		case 2:
			if player1.Mana >= 5 {
				var spell int
				fmt.Println("2 spells available :\n1 - Combo Attack\n2 - Heal")
				fmt.Scanf("%d", &spell)
				switch spell {
				case 1:
					fmt.Println(player1.Name, " has used combo attack !")
					player1.hit(player1, player1.Power)
					player1.hit(player1, player1.Power)
					player1.Mana = 0
					ennemyFightBack(player1, player2)
				case 2:
					fmt.Println(player1.Name, " has used heal !")
					player1.LifeTotal += rand.Intn(10)
					player1.Mana = 0
					ennemyFightBack(player1, player2)
				}
			} else {
				fmt.Println("Not enough Mana going back to action picking..")
			}
		case 3:
			fmt.Println("quitting fight..")
			f = 4
			break
		}
	}
	return 0
}

func gameMenu() {
	var e int
	for e != 4 {
		playerMap := getPlayersFromCSV()
		fmt.Println("Voici la player map \n", playerMap)
		fmt.Println("Voici le menu du jeu")
		fmt.Println("1 - Create a character\n2 - Show Curent player Collection\n3 - Fight against someone else\n4 - Quit")
		fmt.Println("Enter a number to do something..")
		fmt.Scanf("%d", &e)
		switch e {
		case 1:
			// will be moved into register function
			player := Player{}
			fmt.Println("Welcome to the character creation")
			fmt.Println("Tell us your Name..")
			fmt.Scanf("%s", &player.Name)
			player.LifeTotal = 100
			player.Mana = 0
			player.Experience = 0
			player.Power = rand.Intn(40)
			playerMap[len(playerMap)] = player
			writePlayerMapToCSV(playerMap)
		case 2:
			// dev utility
			fmt.Println(playerMap)
		case 3:
			combatMenu(playerMap)
		case 4:
			break
		case 5:
		case 6:
			// readplayer from json soon
		}
	}
}

func combatMenu(playerMap map[int]Player) {
	var j int
	var f int
	fmt.Println(playerMap)
	fmt.Println("Who are you ?")
	fmt.Scanf("%d", &j)
	player := playerMap[j]
	fmt.Println("Who do you want to fight ?")
	fmt.Scanf("%d", &f)
	if j == f {
		println("You can't fight yourself pick someone else")
	} else {
		ennemy := playerMap[f]
		fmt.Println(ennemy)
		winner := combat(&player, &ennemy, j, f)
		fmt.Println("Winner : ", playerMap[winner])
	}
}

func ennemyFightBack(player *Player, ennemy *Player) {
	if player.Mana >= 5 {
		if player.LifeTotal < 50 {
			fmt.Println(ennemy.Name, " has used heal !")
			ennemy.LifeTotal += rand.Intn(10)
			ennemy.Mana = 0
		} else {
			fmt.Println(ennemy.Name, " has used combo attack !")
			ennemy.hit(player, ennemy.Power)
			ennemy.hit(player, ennemy.Power)
			ennemy.Mana = 0
		}
	} else {
		ennemy.hit(player, ennemy.Power)
	}
}

func authenticate() {
	var j int
	var s1 string
	var s2 string
	var s3 string

	fmt.Println("Hello, please register or login if you're already a member of ZEGAMME\n1 - register\n2 - login")
	fmt.Scanf("%d", j)
	if j == 1 {
		fmt.Println("What will your login be ?")
		fmt.Scanf("%s", &s1)

		fmt.Println("and whats your password ?")
		fmt.Scanf("%s", &s2)
		fmt.Println("please confirm your password")
		fmt.Scanf("%s", &s3)

	}

}

func isLoginTaken(login string) {
	//gameMap := getPlayersFromCSV()
}
