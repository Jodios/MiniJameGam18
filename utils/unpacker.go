package utils

import (
	"bytes"
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
)

type TexturePackerJSONArray struct {
	Frames []TexturePackerFrame `json:"frames,omitempty"`
	Meta   struct {
		App     string `json:"app,omitempty"`
		Version string `json:"version,omitempty"`
		Image   string `json:"image,omitempty"`
		Format  string `json:"format,omitempty"`
		Size    struct {
			W int `json:"w,omitempty"`
			H int `json:"h,omitempty"`
		} `json:"size,omitempty"`
		Scale       string `json:"scale,omitempty"`
		Smartupdate string `json:"smartupdate,omitempty"`
	}
}
type TexturePackerFrame struct {
	Filename string `json:"filename,omitempty"`
	Rotated  bool   `json:"rotated,omitempty"`
	Trimmed  bool   `json:"trimmed,omitempty"`
	Frame    struct {
		X int `json:"x,omitempty"`
		Y int `json:"y,omitempty"`
		W int `json:"w,omitempty"`
		H int `json:"h,omitempty"`
	} `json:"frame"`
	SpriteSourceSize struct {
		X int `json:"x,omitempty"`
		Y int `json:"y,omitempty"`
		W int `json:"w,omitempty"`
		H int `json:"h,omitempty"`
	} `json:"spriteSourceSize"`
	SourceSize struct {
		W int `json:"w,omitempty"`
		H int `json:"h,omitempty"`
	} `json:"sourceSize,omitempty"`
}

type ImageWithFrameDetails struct {
	Image     *ebiten.Image
	FrameData TexturePackerFrame
}

type Unpacker struct{}

func (unpacker *Unpacker) UnpackWithFrameDetails(arrayFile []byte, packedImage []byte) (map[string]ImageWithFrameDetails, error) {
	sprites := make(map[string]ImageWithFrameDetails)
	ebitenSpriteSheetImage, texturePackerJSONArray, err := unpacker.parse(arrayFile, packedImage)
	if err != nil {
		return nil, err
	}
	bounds := ebitenSpriteSheetImage.Bounds()
	for _, s := range texturePackerJSONArray.Frames {
		sprites[s.Filename] = ImageWithFrameDetails{
			Image: ebiten.NewImageFromImage(ebitenSpriteSheetImage.SubImage(image.Rectangle{
				Min: image.Point{
					X: bounds.Min.X + s.Frame.X,
					Y: bounds.Min.Y + s.Frame.Y,
				},
				Max: image.Point{
					X: bounds.Min.X + s.Frame.X + s.Frame.W,
					Y: bounds.Min.Y + s.Frame.Y + s.Frame.H,
				},
			})),
			FrameData: s,
		}
	}
	return sprites, nil
}

func (unpacker *Unpacker) Unpack(arrayFile []byte, packedImage []byte) (map[string]*ebiten.Image, error) {
	sprites := make(map[string]*ebiten.Image)
	ebitenSpriteSheetImage, texturePackerJSONArray, err := unpacker.parse(arrayFile, packedImage)
	if err != nil {
		return nil, err
	}
	bounds := ebitenSpriteSheetImage.Bounds()
	for _, s := range texturePackerJSONArray.Frames {
		sprites[s.Filename] = ebitenSpriteSheetImage.SubImage(image.Rectangle{
			Min: image.Point{
				X: bounds.Min.X + s.Frame.X,
				Y: bounds.Min.Y + s.Frame.Y,
			},
			Max: image.Point{
				X: bounds.Min.X + s.Frame.X + s.Frame.W,
				Y: bounds.Min.Y + s.Frame.Y + s.Frame.H,
			},
		}).(*ebiten.Image)
	}
	return sprites, nil
}

func (unpacker *Unpacker) parse(arrayFile []byte, packedImage []byte) (*ebiten.Image, *TexturePackerJSONArray, error) {
	texturePackerJSONArray := new(TexturePackerJSONArray)
	err := json.Unmarshal(arrayFile, texturePackerJSONArray)
	if err != nil {
		return nil, nil, err
	}

	spriteSheetImage, _, err := image.Decode(bytes.NewReader(packedImage))
	if err != nil {
		return nil, nil, err
	}

	ebitenSpriteSheetImage := ebiten.NewImageFromImage(spriteSheetImage)
	return ebitenSpriteSheetImage, texturePackerJSONArray, nil
}
