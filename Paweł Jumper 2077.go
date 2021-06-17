// Imię i nazwisko: Dariusz Rzeźnik
// Grupa: 3
// Nazwa: Prosta gra zręcznościowa
// Tytuł: Paweł Jumper 2077

// Sensem gry jest osiągnięcie najlepszego wyniku (wyrażonego w sekudach) w grze,
// przeciwko osiągniecia tego będą nam uprzykrzać grę potworki,
// aby je ominąć wystarczy je przeskoczyć (spacja)

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	// Variables for running the game in non-standard resolution
	resolutionX int
	resolutionY int

	// Variables for image importing
	img_run      *ebiten.Image
	img_jump     *ebiten.Image
	img_skeleton *ebiten.Image
	img_city     *ebiten.Image
	img_logo     *ebiten.Image
	img_over     *ebiten.Image

	// Variable for main game font
	game_font_small font.Face
	game_font_big   font.Face
)

const (

	// Constants used to show and animate player image
	frame_runOX     = 0
	frame_runOY     = 0
	frameWidth_run  = 110
	frameHeight_run = 113
	frame_run_Num   = 3

	// Constants used to show player jump
	frameWidth_jump  = 120
	frameHeight_jump = 131

	// Constants used to show and animate enemy image
	frame_skeletonOX     = 0
	frame_skeletonOY     = 0
	frameWidth_skeleton  = 79
	frameHeight_skeleton = 120
	frame_skeleton_Num   = 4

	// Constants used to show background image
	frameWidth_city  = 4086
	frameHeight_city = 790

	// Constants used to show and animate logo image
	frame_logoOX     = 0
	frame_logoOY     = 0
	frameWidth_logo  = 600
	frameHeight_logo = 91
	frame_logo_Num   = 5

	// Constants used to show game over image
	frameWidth_over  = 546
	frameHeight_over = 65
)

// Constants describing key pressed at the moment
const (
	KeyNone = iota
	KeySpace
	KeyEscape
)

// Constants describing the state of the game
const (
	ModeMenu = iota
	ModeGame
	ModeGameOV
)

// Structure with main axes
type Position struct {
	X int
	Y int
}

// Structure with game data
type Game struct {
	key_pressed     int
	game_status     int
	player          Position
	jump_rotation   float32
	falldown        bool
	city_background Position
	enemies         []Position
	enemy_speed     int
	timer           int
	score           int
	bestScore       int
}

// Initialization
func init() {
	// Setting seed
	rand.Seed(time.Now().UnixNano())

	// Getting required images
	var err error
	img_run, _, err = ebitenutil.NewImageFromFile("graphics/Man_run.png")
	if err != nil {
		log.Fatal(err)
	}
	img_jump, _, err = ebitenutil.NewImageFromFile("graphics/Man_jump.png")
	if err != nil {
		log.Fatal(err)
	}
	img_city, _, err = ebitenutil.NewImageFromFile("graphics/city_background.png")
	if err != nil {
		log.Fatal(err)
	}
	img_skeleton, _, err = ebitenutil.NewImageFromFile("graphics/Skeleton.png")
	if err != nil {
		log.Fatal(err)
	}
	img_logo, _, err = ebitenutil.NewImageFromFile("graphics/Game_logo.png")
	if err != nil {
		log.Fatal(err)
	}
	img_over, _, err = ebitenutil.NewImageFromFile("graphics/Game_over.png")
	if err != nil {
		log.Fatal(err)
	}

	// Getting custom font
	f, err := os.Open("font/VerminVibes.ttf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.ParseReaderAt(f)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	game_font_small, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	game_font_big, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// A function that resets game data
func (g *Game) reset() {
	g.key_pressed = KeyNone
	g.game_status = ModeMenu

	g.score = 0
	for i := 0; i < 5; i++ {
		g.enemies[i].X = (resolutionX + 50) + (i * 500) + (rand.Intn(60) - 30)
		g.enemies[i].Y = (resolutionY / 4) * 3
	}
	g.enemy_speed = 1
	g.player.X = (resolutionX / 5)
	g.player.Y = (resolutionY / 4) * 3
	g.falldown = false
	g.jump_rotation = 0.0
	g.city_background.X = 0
	g.city_background.Y = (resolutionY / 4) * 3
}

// A function that controls the player's position while jumping
func (g *Game) Jump() {
	if !g.falldown {
		g.player.Y -= 2 * g.enemy_speed
		g.jump_rotation -= 0.2
		if g.player.Y <= (resolutionY/4)*3-170 {
			g.player.Y = (resolutionY/4)*3 - 170
			g.falldown = true
		}
	} else if g.falldown {
		g.player.Y += 2 * g.enemy_speed
		g.jump_rotation += 0.2
		if g.player.Y >= (resolutionY/4)*3 {
			g.player.Y = (resolutionY / 4) * 3
			g.falldown = false
			g.key_pressed = KeyNone
		}
	}
}

// A function that returns if player hits the enemy or not
func (g *Game) Colision() bool {

	x0 := g.player.X - (frameWidth_run / 2)
	x1 := x0 + frameWidth_run - 10
	y1 := g.player.Y
	for i := 0; i < 5; i++ {
		if x1 >= g.enemies[i].X && x0+25 <= g.enemies[i].X && y1 >= g.enemies[i].Y-100 {
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	// Game's logical update.

	switch g.game_status {
	case ModeMenu:
		if g.key_pressed == KeyEscape {
			fmt.Printf("EXIT\n")
			os.Exit(0)
		} else if g.key_pressed == KeySpace {
			g.game_status = ModeGame
		}
	case ModeGame:

		if g.Colision() {
			g.game_status = ModeGameOV
			g.key_pressed = KeyNone
		}

		// Updating score and game speed
		if g.timer%60 == 0 {
			g.score++
			if g.score == 35 {
				g.enemy_speed++
			}
			if g.score == 55 {
				g.enemy_speed++
			}
		}

		if g.key_pressed == KeyEscape {
			g.reset()
		}

		if g.key_pressed == KeySpace {
			g.Jump()
		}

		// Updating enemies speed and position
		for i := 0; i < 4; i++ {
			g.enemies[i].X -= 2 * g.enemy_speed
			if g.enemies[i].X <= 0 {
				g.enemies[i].X = g.enemies[4].X - 800 + (rand.Intn(450))
			}
		}

		// Updating background animation
		g.city_background.X -= 1
		if g.city_background.X == -682 {
			g.city_background.X = 0
		}

		if g.score >= g.bestScore {
			g.bestScore = g.score
		}

	case ModeGameOV:
		if g.key_pressed == KeySpace {
			g.reset()
		}
	}

	g.timer++

	// Getting key value from keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.key_pressed = KeySpace
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.key_pressed = KeyEscape
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Game's rendering.

	if g.game_status == ModeMenu {

		// Drawing game logo
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(frameWidth_logo/2), -float64(frameHeight_logo/2))
		op.GeoM.Translate(float64(resolutionX/2), float64(resolutionY/2))
		i := (g.timer / (rand.Intn(5000) + 1)) % frame_logo_Num
		sx, sy := frame_logoOX+i*frameWidth_logo, frame_logoOY
		screen.DrawImage(img_logo.SubImage(image.Rect(sx, sy, sx+frameWidth_logo, sy+frameHeight_logo)).(*ebiten.Image), op)

		// Drawing menu options
		text.Draw(screen, fmt.Sprintf("Press SPACE to start"), game_font_small, resolutionX/4, resolutionY/4*3, color.White)
		text.Draw(screen, fmt.Sprintf("Press ESCAPE to quit"), game_font_small, resolutionX/4, resolutionY/4*3+50, color.White)

	} else if g.game_status == ModeGame {

		// Drawing keybinds hints
		text.Draw(screen, fmt.Sprintf("Escape: RESTART        Space: JUMP "), game_font_small, 60, resolutionY-60, color.White)

		// Drawing background
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(0), -float64(frameHeight_city)-10)
		op.GeoM.Translate(float64(g.city_background.X), float64(g.city_background.Y))
		screen.DrawImage(img_city, op)

		// Drawing score value
		text.Draw(screen, fmt.Sprintf("Score %d", g.score), game_font_big, resolutionX/4*3, resolutionY/7, color.White)

		// Drawing player hit-box
		// ebitenutil.DrawRect(screen, float64(g.player.X-(frameWidth_run/2)+35), float64(g.player.Y-frameHeight_run), float64(frameWidth_run-45), float64(frameHeight_run), color.RGBA{0xff, 0xff, 0xff, 0xff})

		//Drawing enemies
		for i := 0; i < 5; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-float64(frameWidth_skeleton/2), -float64(frameHeight_skeleton))
			op.GeoM.Translate(float64(g.enemies[i].X), float64(g.enemies[i].Y))
			i := (g.timer / 15) % frame_skeleton_Num
			sx, sy := frame_skeletonOX+i*frameWidth_skeleton, frame_skeletonOY
			screen.DrawImage(img_skeleton.SubImage(image.Rect(sx, sy, sx+frameWidth_skeleton, sy+frameHeight_skeleton)).(*ebiten.Image), op)
		}

		// Drawing player while jumping
		if g.key_pressed == KeySpace {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-float64(frameWidth_jump)/2, -float64(frameHeight_jump)+10)
			op.GeoM.Rotate(math.Pi * float64(g.jump_rotation) / 250)
			op.GeoM.Translate(float64(g.player.X), float64(g.player.Y))
			screen.DrawImage(img_jump, op)
		} else { // Drawing player while running
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-float64(frameWidth_run)/2, -float64(frameHeight_run))
			op.GeoM.Translate(float64(g.player.X), float64(g.player.Y))
			i := (g.timer / 10) % frame_run_Num
			sx, sy := frame_runOX+i*frameWidth_run, frame_runOY
			screen.DrawImage(img_run.SubImage(image.Rect(sx, sy, sx+frameWidth_run, sy+frameHeight_run)).(*ebiten.Image), op)
		}

	} else if g.game_status == ModeGameOV {
		// Drawing "game over" image
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(frameWidth_over/2), -float64(frameHeight_over))
		op.GeoM.Translate(float64(resolutionX/2), float64(resolutionY/3))
		screen.DrawImage(img_over, op)

		// Drawig options
		if g.score == g.bestScore {
			text.Draw(screen, fmt.Sprintf("NEW BEST SCORE: %d", g.bestScore), game_font_big, resolutionX/4, resolutionY/2, color.White)
			text.Draw(screen, fmt.Sprintf("Press SPACE to continue"), game_font_big, resolutionX/4, resolutionY/2+50, color.White)
		} else {
			text.Draw(screen, fmt.Sprintf("SCORE: %d", g.score), game_font_big, resolutionX/4, resolutionY/2, color.White)
			text.Draw(screen, fmt.Sprintf("BEST SCORE: %d", g.bestScore), game_font_big, resolutionX/4, resolutionY/2+50, color.White)
			text.Draw(screen, fmt.Sprintf("Press SPACE to continue"), game_font_big, resolutionX/4, resolutionY/2+100, color.White)
		}
	}

	ebitenutil.DrawRect(screen, float64(0), float64(0), float64(resolutionX), 30, color.RGBA{0x00, 0x00, 0x00, 0xff}) //Top Black
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()), 30, 5)

	// DRAWING TEST LINES
	// ebitenutil.DrawLine(screen, 0, float64(resolutionY)/2, float64(resolutionX), float64(resolutionY)/2, color.RGBA{0xff, 0xff, 0xff, 0xff})
	// ebitenutil.DrawLine(screen, float64(g.player.X), 0, float64(g.player.X), float64(resolutionY), color.RGBA{0xff, 0xff, 0xff, 0xff})

	// Drawing window frame borders
	ebitenutil.DrawRect(screen, float64(0), float64(resolutionY-50), float64(resolutionX), 20, color.RGBA{0xff, 0xff, 0xff, 0xff}) //Bottom
	ebitenutil.DrawRect(screen, float64(0), float64(30), float64(resolutionX), 20, color.RGBA{0xff, 0xff, 0xff, 0xff})             //Top
	ebitenutil.DrawRect(screen, float64(0), float64(0), 20, float64(resolutionY), color.RGBA{0xff, 0xff, 0xff, 0xff})              //Left
	ebitenutil.DrawRect(screen, float64(resolutionX-20), float64(0), 20, float64(resolutionY), color.RGBA{0xff, 0xff, 0xff, 0xff}) //Right
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return resolutionX, resolutionY
}

func newGame() *Game {
	g := &Game{
		enemies: make([]Position, 5),
	}

	g.key_pressed = KeyNone
	g.game_status = ModeMenu

	g.score = 0
	g.bestScore = 0
	for i := 0; i < 5; i++ {
		g.enemies[i].X = (resolutionX + 50) + (i * 500) + (rand.Intn(60) - 30)
		g.enemies[i].Y = (resolutionY / 4) * 3
	}
	g.enemy_speed = 1
	g.player.X = (resolutionX / 5)
	g.player.Y = (resolutionY / 4) * 3
	g.falldown = false
	g.jump_rotation = 0.0
	g.city_background.X = 0
	g.city_background.Y = (resolutionY / 4) * 3

	return g
}

func main() {

	// Flags for non-default resolutions
	flag.IntVar(&resolutionX, "resX", 640, "Rozdzielczość \"X\"")
	flag.IntVar(&resolutionY, "resY", 480, "Rozdzielczość \"Y\"")
	flag.Parse()

	// Resolution min-size check
	if resolutionX < 640 || resolutionY < 480 {
		resolutionX = 640
		resolutionY = 480
		fmt.Printf("Too small resolution\nChanged to default\n")
	}
	fmt.Printf("Game launched with %d:%d resolution\n", resolutionX, resolutionY)

	// Setting window resolution and title
	ebiten.SetWindowSize(resolutionX, resolutionY)
	ebiten.SetWindowTitle("PAWEŁ JUMPER 2077")

	// Starting the game
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
