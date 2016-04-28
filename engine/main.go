package main

import (
	"image/color"
	"log"
	"time"

	"github.com/ttfx-bordeaux/code-of-war-5/engine/game"

	"engo.io/ecs"
	"engo.io/engo"
)

var (
	screenWidth                 = 1024
	screenHeight                = 640
	TileWidth       int         = 120
	TileHeight      int         = 120
	nbTilesAbs      int         = 3
	nbTilesOrd      int         = 5
	ImgGroundName   string      = "grass-600-600.png"
	ImgTowerName    string      = "tour1-600-600.png"
	backgroundColor color.Color = color.White
	padRight        int
	players         map[string]*Player
)

func toPoint(p game.GamePosition, loop int) engo.Point {
	padding := loop * (nbTilesAbs*TileWidth + padRight)
	return engo.Point{X: float32(p.Abs*TileWidth + padding), Y: float32(p.Ord * TileHeight)}
}

type Player struct {
	Id       string
	Ground   []*ecs.Entity
	towers   []*ecs.Entity
	chickens []*ecs.Entity
}

//GameWorld world
type GameWorld struct{}

//Preload preload
func (world *GameWorld) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engo.Files.AddFromDir("static", false)

	game.Preload()

	log.Println("Preloaded")
}

//Setup setup
func (world *GameWorld) Setup(w *ecs.World) {
	engo.SetBackground(backgroundColor)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.AnimationSystem{})
	w.AddSystem(&game.ControlSystem{})

	w.AddSystem(&engo.AudioSystem{})
	w.AddSystem(&game.WhoopSystem{})

	createBgMusic(w)

	idGround := 0
	var chicken *ecs.Entity
	for _, p := range players {
		log.Println(p)
		log.Println(idGround)
		createGround(w, idGround, ImgGroundName)
		createTower(w, game.GamePosition{1, 3}, idGround, ImgTowerName)
		chicken = createChicken(w, game.GamePosition{0, 0}, idGround)
		idGround++
	}
	go loopGame(chicken, 1)
}

func loopGame(chicken *ecs.Entity, idGround int) {
	// loop game
	//var time *engo.Clock
	ticker := time.NewTicker(time.Duration(int(time.Second)))
Outer:
	for {
		select {
		case <-ticker.C:
			cf := &game.GamePosition{}
			space := &engo.SpaceComponent{}
			var ok bool
			if cf, ok = chicken.ComponentFast(cf).(*game.GamePosition); !ok {
				log.Println("no gamePosition")
				break Outer
			}
			if cf.Ord < nbTilesOrd {
				if space, ok = chicken.ComponentFast(space).(*engo.SpaceComponent); !ok {
					log.Println("no spaceComponent")
					break Outer
				}
				cf.Ord++
				space.Position = toPoint(*cf, idGround)
				log.Printf("move")
			} else {
				break Outer
			}
			// call http
			//if close {
			//	break Outer
			// }
			// if !headless && window.ShouldClose() {
			// closeEvent()
			// }
		}
	}
	log.Println("stop")
	ticker.Stop()
}

func createBgMusic(w *ecs.World) {
	bgMusic := ecs.NewEntity("AudioSystem", "WhoopSystem")
	bgMusic.AddComponent(&engo.AudioComponent{File: "sound.wav", Repeat: true, Background: true, RawVolume: 1})
	if err := w.AddEntity(bgMusic); err != nil {
		log.Println(err)
	}
}

func createTower(w *ecs.World, p game.GamePosition, idGround int, imgName string) {
	tower := createEntityTower(imgName, toPoint(p, idGround))
	err := w.AddEntity(tower)
	if err != nil {
		log.Println(err)
	}
}

func createGround(w *ecs.World, idGround int, imgName string) {
	padding := idGround * (nbTilesAbs*TileWidth + padRight)
	for j := 0; j < nbTilesOrd; j++ {
		for i := 0; i < nbTilesAbs; i++ {
			grass := createEntityTile(imgName, engo.Point{X: float32(i*TileWidth + padding), Y: float32(j * TileHeight)})
			if err := w.AddEntity(grass); err != nil {
				log.Println(err)
			}
		}
	}
}

func createChicken(w *ecs.World, p game.GamePosition, idGround int) *ecs.Entity {
	spriteSheet := engo.NewSpritesheetFromFile("chicken.png", 150, 150)
	entity := createEntityChicken(toPoint(p, idGround), spriteSheet, game.StopAction)

	entity.AddComponent(&p)

	err := w.AddEntity(entity)
	if err != nil {
		log.Println(err)
	}

	return entity
}

func createEntityChicken(point engo.Point, spriteSheet *engo.Spritesheet, action *engo.AnimationAction) *ecs.Entity {
	entity := ecs.NewEntity("AnimationSystem", "RenderSystem", "ControlSystem")

	render := engo.NewRenderComponent(spriteSheet.Cell(action.Frames[0]), engo.Point{X: 1, Y: 1}, "chicken")
	width := 150 * render.Scale().X
	height := 150 * render.Scale().Y
	space := &engo.SpaceComponent{Position: point, Width: width, Height: height}
	animation := engo.NewAnimationComponent(spriteSheet.Drawables(), game.AnimationRate)
	animation.AddAnimationActions(game.Actions)
	animation.SelectAnimationByAction(action)

	entity.AddComponent(render)
	entity.AddComponent(space)
	entity.AddComponent(animation)

	return entity
}

func createEntityTile(imgName string, point engo.Point) *ecs.Entity {
	entity := ecs.NewEntity("RenderSystem")

	texture := engo.Files.Image(imgName)
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}

	render := engo.NewRenderComponent(texture, engo.Point{X: 0.2, Y: 0.2}, "tile")
	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y
	space := &engo.SpaceComponent{Position: point, Width: width, Height: height}

	entity.AddComponent(render)
	entity.AddComponent(space)

	return entity
}

func createEntityTower(imgName string, point engo.Point) *ecs.Entity {
	entity := ecs.NewEntity("RenderSystem")

	texture := engo.Files.Image(imgName)
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}
	render := engo.NewRenderComponent(texture, engo.Point{X: 0.2, Y: 0.2}, "tile")
	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y
	space := &engo.SpaceComponent{Position: point, Width: width, Height: height}

	entity.AddComponent(render)
	entity.AddComponent(space)

	return entity
}

//Hide hide
func (*GameWorld) Hide() {}

//Show show
func (*GameWorld) Show() {}

//Exit
func (*GameWorld) Exit() {
	log.Println("[GAME] Exit event called")
	//Here if you want you can prompt the user if they're sure they want to close
	engo.Exit()
}

//Type type
func (*GameWorld) Type() string { return "GameWorld" }

func main() {
	// see level.go
	// tick, fps
	players = make(map[string]*Player, 2)
	players["1"] = &Player{Id: "1"}
	players["2"] = &Player{Id: "2"}
	sizeWidth := TileWidth * nbTilesAbs * len(players)
	if sizeWidth > screenWidth {
		padRight = 10
	} else {
		padRight = (screenWidth - sizeWidth) / len(players)
	}
	opts := engo.RunOptions{
		Title:      "Code of War : Enlarge your tower",
		Width:      screenWidth,
		Height:     screenHeight,
		Fullscreen: false,
	}
	engo.OverrideCloseAction()
	engo.Run(opts, &GameWorld{})
}
