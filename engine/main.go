package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
)

var (
	screenWidth                 = 1024
	screenHeight                = 750
	TileWidth       int         = 32
	TileHeight      int         = 32
	nbTilesAbs      int         = 15
	nbTilesOrd      int         = 20
	ImgGroundName   string      = "tile-desert-32.png"
	ImgTowerName    string      = "tour1-600-600.png"
	backgroundColor color.Color = color.White
	padRight        int
	players         map[string]*Player
)

type Player struct {
	Id       string
	IdGround int
	ground   []*tileEntity
	towers   []*towerEntity
	chickens []*chickenEntity
}

type towerEntity struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

type chickenEntity struct {
	ecs.BasicEntity
	engo.AnimationComponent
	engo.RenderComponent
	engo.SpaceComponent
	PositionComponent
}

type tileEntity struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

//GameWorld world
type DefaultScene struct{}

//Preload preload
func (world *DefaultScene) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engo.Files.AddFromDir("static/img", false)

	Preload()

	log.Println("Preloaded")
}

//Setup setup
func (scene *DefaultScene) Setup(w *ecs.World) {
	log.Println("call setup")
	engo.SetBackground(backgroundColor)

	renderSystem := &engo.RenderSystem{}
	animationSystem := &engo.AnimationSystem{}
	controlSystem := &ControlSystem{}
	moveSystem := &MoveSystem{}

	w.AddSystem(controlSystem)
	w.AddSystem(moveSystem)
	w.AddSystem(animationSystem)
	w.AddSystem(renderSystem)

	for _, player := range players {
		ground := createGround(w, player.IdGround, ImgGroundName, renderSystem)
		tower := createTower(w, NewPositionComponent(1, 3, player.IdGround), ImgTowerName, renderSystem)
		chicken := createChicken(w, NewPositionComponent(0, 0, player.IdGround), renderSystem, animationSystem, controlSystem, moveSystem)

		player.ground = ground
		player.towers = append(player.towers, tower)
		player.chickens = append(player.chickens, chicken)
	}

	players["1"].chickens[0].PositionComponent.changePositionTo(14, 19)
	players["2"].chickens[0].PositionComponent.changePositionTo(10, 2)
}

func createTower(w *ecs.World, p PositionComponent, imgName string, renderSystem *engo.RenderSystem) *towerEntity {
	tower := createEntityTower(imgName, p.toPoint())
	renderSystem.Add(&tower.BasicEntity, &tower.RenderComponent, &tower.SpaceComponent)
	return tower
}

func createGround(w *ecs.World, idGround int, imgName string, renderSystem *engo.RenderSystem) []*tileEntity {
	ground := make([]*tileEntity, nbTilesOrd*nbTilesAbs)
	padding := idGround * (nbTilesAbs*TileWidth + padRight)
	for j := 0; j < nbTilesOrd; j++ {
		for i := 0; i < nbTilesAbs; i++ {
			tile := createEntityTile(imgName, engo.Point{X: float32(i*TileWidth + padding), Y: float32(j * TileHeight)})
			renderSystem.Add(&tile.BasicEntity, &tile.RenderComponent, &tile.SpaceComponent)
			ground = append(ground, tile)
		}
	}
	return ground
}

func createChicken(w *ecs.World, p PositionComponent,
	renderSystem *engo.RenderSystem, animationSystem *engo.AnimationSystem,
	controlSystem *ControlSystem, moveSystem *MoveSystem) *chickenEntity {

	entity := createEntityChicken(p)

	renderSystem.Add(&entity.BasicEntity, &entity.RenderComponent, &entity.SpaceComponent)
	//animationSystem.Add(&entity.BasicEntity, &entity.AnimationComponent, &entity.RenderComponent)
	//controlSystem.Add(&entity.BasicEntity, &entity.AnimationComponent)
	moveSystem.Add(&entity.BasicEntity, &entity.PositionComponent, &entity.SpaceComponent)

	return entity;
}

func createEntityChicken(p PositionComponent) *chickenEntity {
	spriteSheet := engo.NewSpritesheetFromFile("chicken.png", 150, 150)

	entity := &chickenEntity{BasicEntity: ecs.NewBasic()}

	texture := engo.Files.Image("chicken-32.png")

	//entity.SpaceComponent = engo.SpaceComponent{Position: point, Width: 150, Height: 150}
	//entity.RenderComponent = engo.NewRenderComponent(spriteSheet.Cell(0), engo.Point{X: 1, Y: 1})
	entity.SpaceComponent = engo.SpaceComponent{Position: p.toPoint(), Width: texture.Width(), Height: texture.Height()}
	entity.RenderComponent = engo.NewRenderComponent(texture, engo.Point{X: 1, Y: 1})
	entity.AnimationComponent = engo.NewAnimationComponent(spriteSheet.Drawables(), AnimationRate)
	entity.AnimationComponent.AddAnimations(Actions)
	entity.AnimationComponent.AddDefaultAnimation(StopAction)
	entity.PositionComponent = p

	return entity
}

func createEntityTile(imgName string, point engo.Point) *tileEntity {
	entity := &tileEntity{BasicEntity: ecs.NewBasic()}

	texture := engo.Files.Image(imgName)

	entity.RenderComponent = engo.NewRenderComponent(texture, engo.Point{X: 1, Y: 1})
	entity.SpaceComponent = engo.SpaceComponent{Position: point, Width: texture.Width(), Height: texture.Height()}

	return entity
}

func createEntityTower(imgName string, point engo.Point) *towerEntity {
	entity := &towerEntity{BasicEntity: ecs.NewBasic()}

	texture := engo.Files.Image(imgName)
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}
	entity.SpaceComponent = engo.SpaceComponent{Position: point, Width: texture.Width(), Height: texture.Height()}
	entity.RenderComponent = engo.NewRenderComponent(texture, engo.Point{X: 0.2, Y: 0.2})

	return entity
}

//Hide hide
func (*DefaultScene) Hide() {}

//Show show
func (*DefaultScene) Show() {}

//Type type
func (*DefaultScene) Type() string { return "GameWorld" }

func main() {
	// see level.go
	// tick, fps
	players = make(map[string]*Player, 2)
	players["1"] = &Player{Id: "1", IdGround:0}
	players["2"] = &Player{Id: "2", IdGround:1}
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
	engo.Run(opts, &DefaultScene{})
}
