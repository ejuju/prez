package prez

import (
	_ "embed"
	"fmt"
	"io"
	"time"

	"github.com/go-pdf/fpdf"
)

const (
	marginSide       = 50
	lineHeight       = fontSize + 6 // Only for regular lines.
	pageHeight       = 595.28
	topMargin        = 50
	fontSizeTitle    = 35
	fontSizeSubtitle = 20
	fontSizeH1       = 20
	fontSize         = 10
	fontFamily       = "sans-serif"
	fontFamilyMono   = "monospace"
)

var (
	//go:embed noto-sans-regular.ttf
	notoRegularTTF []byte
	//go:embed noto-sans-bold.ttf
	notoBoldTTF []byte
	//go:embed jetbrains-mono-regular.ttf
	jetbrainsMonoRegularTTF []byte
)

func WritePDF(w io.Writer, doc *Document) (err error) {
	// Init PDF and add basic metadata.
	pdf := fpdf.New("L", "pt", "A4", "")
	pdf.SetCreationDate(time.Now())
	pdf.SetAuthor(doc.Author, true)
	pdf.SetTitle(doc.Title, true)
	pdf.SetLang(doc.Lang)

	// Set font.
	pdf.AddUTF8FontFromBytes(fontFamily, "", notoRegularTTF)
	pdf.AddUTF8FontFromBytes(fontFamily, "B", notoBoldTTF)
	pdf.AddUTF8FontFromBytes(fontFamilyMono, "", jetbrainsMonoRegularTTF)
	pdf.SetFont(fontFamily, "", fontSize)

	// Set global page margins (and init color).
	pdf.SetTopMargin(topMargin)
	pdf.SetLeftMargin(marginSide)
	pdf.SetRightMargin(marginSide)
	pdf.SetTextColor(10, 10, 10)

	// Render cover (= first page).
	pdf.AddPage()
	pdf.SetFontSize(fontSizeTitle)
	pdf.SetFontStyle("B")
	pdf.MultiCell(0, fontSizeTitle, doc.Title, "", "C", false)
	pdf.Ln(fontSizeSubtitle)
	pdf.SetFontSize(fontSizeSubtitle)
	pdf.SetFontStyle("")
	pdf.MultiCell(0, fontSizeSubtitle, doc.Author, "", "C", false)

	// Render pages.
	for _, page := range doc.Pages {
		pdf.AddPage()
		for _, block := range page {
			switch block := block.(type) {
			default:
				return fmt.Errorf("unsupported block type: %T", block)
			case H1:
				writeH1(pdf, string(block))
			case Text:
				writeTextLine(pdf, string(block))
			case Image:
				writeImage(pdf, string(block))
			case Code:
				writeCode(pdf, string(block))
			case ListItem:
				writeListItem(pdf, string(block))
			}
		}
	}

	err = pdf.Output(w)
	if err != nil {
		return fmt.Errorf("render PDF: %w", err)
	}
	return nil
}

func writeH1(pdf *fpdf.Fpdf, v string) {
	pdf.SetFontSize(fontSizeH1)
	pdf.SetFontStyle("B")
	pdf.MultiCell(0, fontSizeH1, v, "", "L", false)
	pdf.Ln(lineHeight)
}

func writeTextLine(pdf *fpdf.Fpdf, v string) {
	pdf.SetFontSize(fontSize)
	pdf.SetFontStyle("")
	pdf.MultiCell(0, fontSize+4, v, "", "L", false)
}

func writeCode(pdf *fpdf.Fpdf, v string) {
	pdf.Ln(fontSize)
	pdf.SetFont(fontFamilyMono, "", fontSize)
	pdf.SetFontStyle("")
	pdf.MultiCell(0, fontSize+4, v, "1", "L", false)
	pdf.SetFont(fontFamily, "", fontSize) // Reset default font.
}

func writeListItem(pdf *fpdf.Fpdf, v string) {
	writeTextLine(pdf, "- "+v)
}

func writeImage(pdf *fpdf.Fpdf, fpath string) {
	pdf.ImageOptions(fpath, marginSide, 0, 0, 0, true, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
}
