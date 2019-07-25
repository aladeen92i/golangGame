package main

import (
	"encoding/csv"
	"fmt"
	"io"
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

func writePlayerMapToCSV(playerMap map[int]playerStruct) {
	file, err := os.Create("players.csv")
	checkError("Error:", err)
	defer file.Close()
	headers := []string{
		"name",
		"lifeTotal",
		"power",
		"mana",
		"experience",
	}
	// write column headers
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for key := range playerMap {
		r := make([]string, 0, 1+len(headers))
		lifeString := strconv.Itoa(playerMap[key].lifeTotal)
		powerString := strconv.Itoa(playerMap[key].power)
		manaString := strconv.Itoa(playerMap[key].mana)
		experienceString := strconv.Itoa(playerMap[key].experience)
		r = append(
			r,
			playerMap[key].name,
			lifeString,
			powerString,
			manaString,
			experienceString,
		)
		writer.Write(r)
	}
}

func getPlayersFromCSV() map[int]playerStruct {
	var cpt int
	playerMap := make(map[int]playerStruct)
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
		lifeTotal, err := strconv.Atoi(record[1])
		power, err := strconv.Atoi(record[2])
		mana, err := strconv.Atoi(record[3])
		experience, err := strconv.Atoi(record[4])
		player := playerStruct{record[0], lifeTotal, power, mana, experience}
		playerMap[cpt] = player
		cpt += 1
		//fmt.Printf("name: %s lifeTotal: %s power: %s mana: %s experience: %s\n", record[0], record[1], record[2], record[3], record[4])
	}
	return playerMap
}

type playerStruct struct {
	name       string
	lifeTotal  int
	power      int
	mana       int
	experience int
}

func (p *playerStruct) hit(playerHit *playerStruct, power int) {
	fmt.Println(p.name, "tries to hit ", playerHit.name, " with a power of ", power)
	playerHit.getHit(power)
	p.mana += 1
	//fmt.Println(p.name, " has ", p.mana, " mana points ")
}

func (p *playerStruct) getHit(power int) {
	p.lifeTotal = p.lifeTotal - power
	//fmt.Println(p.name, " has took a torgnole and now has ", p.lifeTotal, " HP ")
}

func (p *playerStruct) getXp() {
	p.experience += 1
	//fmt.Println(p.name, " has took a torgnole and now has ", p.lifeTotal, " HP ")
}

func combat(player1 *playerStruct, player2 *playerStruct, key1 int, key2 int) int {
	var f int
	for f != 4 {
		if player1.lifeTotal <= 0 {
			fmt.Println(player2)
			fmt.Println(player2.name, " has won !")
			player1.lifeTotal = 100
			player2.lifeTotal = 100
			player1.mana = 0
			player2.mana = 0
			player2.experience = player2.experience + 1
			f = 4
			return key2
		}
		if player2.lifeTotal <= 0 {
			fmt.Println(player1)
			fmt.Println(player1.name, " has won !")
			player1.lifeTotal = 100
			player2.lifeTotal = 100
			player1.mana = 0
			player2.mana = 0
			player1.experience = player1.experience + 1
			f = 4
			return key1
		}
		fmt.Println("You're in a fight, what do you do ?")
		fmt.Println("Your life: ", player1.lifeTotal, ", your mana : ", player1.mana)
		fmt.Println("Your ennemy : ", player2.name, " ", player2.lifeTotal, " HP , ", player2.mana, " mana points ")
		fmt.Println("1 - Simple Attack\n2 - Spell\n3 - run away ??")
		fmt.Scanf("%d", &f)
		switch f {
		case 1:
			player1.hit(player2, player1.power)
			ennemyFightBack(player1, player2)
		case 2:
			if player1.mana >= 5 {
				var spell int
				fmt.Println("2 spells available :\n1 - Combo Attack\n2 - Heal")
				fmt.Scanf("%d", &spell)
				switch spell {
				case 1:
					fmt.Println(player1.name, " has used combo attack !")
					player1.hit(player1, player1.power)
					player1.hit(player1, player1.power)
					player1.mana = 0
					ennemyFightBack(player1, player2)
				case 2:
					fmt.Println(player1.name, " has used heal !")
					player1.lifeTotal += rand.Intn(10)
					player1.mana = 0
					ennemyFightBack(player1, player2)
				}
			} else {
				fmt.Println("Not enough mana going back to action picking..")
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
			player := playerStruct{}
			fmt.Println("Welcome to the character creation")
			fmt.Println("Tell us your name..")
			fmt.Scanf("%s", &player.name)
			player.lifeTotal = 100
			player.mana = 0
			player.experience = 0
			player.power = rand.Intn(40)
			playerMap[len(playerMap)] = player
			writePlayerMapToCSV(playerMap)
		case 2:
			// dev utility
			fmt.Println(playerMap)
		case 3:
			combatMenu(playerMap)
		case 4:
			break
		}
	}
}

func combatMenu(playerMap map[int]playerStruct) {
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

func ennemyFightBack(player *playerStruct, ennemy *playerStruct) {
	if player.mana >= 5 {
		if player.lifeTotal < 50 {
			fmt.Println(ennemy.name, " has used heal !")
			ennemy.lifeTotal += rand.Intn(10)
			ennemy.mana = 0
		} else {
			fmt.Println(ennemy.name, " has used combo attack !")
			ennemy.hit(player, ennemy.power)
			ennemy.hit(player, ennemy.power)
			ennemy.mana = 0
		}
	} else {
		ennemy.hit(player, ennemy.power)
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
		fmt.Scanf("%s", &s1)
		fmt.Println("please confirm your password")
		fmt.Scanf("%s", &s1)

	}

}

func isLoginTaken(login string) {

}
