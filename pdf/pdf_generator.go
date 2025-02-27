package pdf

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type Task struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderID   uuid.UUID
	Name      string
	Data      string
	Status    string
}

func GenerateDocument(data []Task) ([]byte, error) {
	cfg := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()

	pdfFile := maroto.New(cfg)

	//addHeader(pdfFile)

	addDetails(pdfFile, len(data))

	addItemList(pdfFile, data)

	doc, errGen := pdfFile.Generate()
	if errGen != nil {
		return nil, errGen
	}

	return doc.GetBytes(), nil
}

func addHeader(m core.Maroto) {
	m.AddRow(20,
		text.NewCol(12, "Header",
			props.Text{
				Top:   5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  16,
			}))
	m.AddRow(20,
		text.NewCol(12, "ORDERS",
			props.Text{
				Top:   5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  16,
			}))

}

func addDetails(m core.Maroto, itemcount int) {
	m.AddRow(10,
		text.NewCol(6, "Date: "+time.Now().Format(time.RFC822),
			props.Text{
				Align: align.Left,
				Size:  10,
			}),
		text.NewCol(6, fmt.Sprint("Total orders: ", itemcount), props.Text{
			Align: align.Right,
			Size:  10,
		}))
	m.AddRow(10, line.NewCol(12))
}

func (t Task) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(2, "ID", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Created at", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Updated at", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Order id", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Name", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Data length", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Status", props.Text{Style: fontstyle.Bold}),
	)
}

func (t Task) GetContent(i int) core.Row {
	r := row.New(10).Add(
		text.NewCol(2, "ID", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Order id", props.Text{Style: fontstyle.Bold})).
		Add(
			text.NewCol(5, t.ID.String()),
			text.NewCol(5, t.OrderID.String())).
		Add(
			text.NewCol(2, "Created at", props.Text{Style: fontstyle.Bold}),
			text.NewCol(2, "Updated at", props.Text{Style: fontstyle.Bold}),
			text.NewCol(2, "Name", props.Text{Style: fontstyle.Bold}),
			text.NewCol(2, "Data length", props.Text{Style: fontstyle.Bold}),
			text.NewCol(2, "Status", props.Text{Style: fontstyle.Bold})).
		Add(
			text.NewCol(2, t.CreatedAt.Format(time.RFC822)),
			text.NewCol(2, t.UpdatedAt.Format(time.RFC822)),
			text.NewCol(2, t.Name),
			text.NewCol(2, fmt.Sprint(len(t.Data))),
			text.NewCol(2, t.Status))

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
		})
	}
	return r
}

func addItemList(m core.Maroto, data []Task) error {
	rows, err := list.Build(data)
	if err != nil {
		return err
	}
	m.AddRows(rows...)
	return nil
}
