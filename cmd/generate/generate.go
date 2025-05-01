package generate

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/spf13/cobra"

	"github.com/maaslalani/invoice/inter"
)

type invoice struct {
	Id    string `json:"id"    yaml:"id"`
	Title string `json:"title" yaml:"title"`

	Logo string `json:"logo" yaml:"logo"`
	From string `json:"from" yaml:"from"`
	To   string `json:"to"   yaml:"to"`
	Date string `json:"date" yaml:"date"`
	Due  string `json:"due"  yaml:"due"`

	Items      []string  `json:"items"      yaml:"items"`
	Quantities []float64 `json:"quantities" yaml:"quantities"`
	Rates      []float64 `json:"rates"      yaml:"rates"`

	Tax      float64 `json:"tax"      yaml:"tax"`
	Discount float64 `json:"discount" yaml:"discount"`
	Currency string  `json:"currency" yaml:"currency"`

	Note string `json:"note" yaml:"note"`

	ForeignExchange foreignExchange

	Output string
}

type foreignExchange struct {
	Rate     float64 `json:"rate"     yaml:"rate"`
	Currency string  `json:"currency" yaml:"currency"`
}

func getDefaultInvoice() invoice {
	return invoice{
		Id:         time.Now().Format("20060102"),
		Title:      "INVOICE",
		Logo:       "",
		From:       "Project Folded, Inc.",
		To:         "Untitled Corporation, Inc.",
		Date:       time.Now().Format("Jan 02, 2006"),
		Due:        time.Now().AddDate(0, 0, 14).Format("Jan 02, 2006"),
		Items:      []string{"Paper Cranes"},
		Quantities: []float64{2},
		Rates:      []float64{25},
		Tax:        0,
		Discount:   0,
		Currency:   "USD",
		Note:       "",
		ForeignExchange: foreignExchange{
			Rate:     1,
			Currency: "USD",
		},
	}
}

func generate(file *invoice) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{
			PageSize: *gopdf.PageSizeA4,
		})

		pdf.SetMargins(40, 40, 40, 40)
		pdf.AddPage()

		err := pdf.AddTTFFontData("Inter", inter.Font)
		if err != nil {
			return err
		}

		err = pdf.AddTTFFontData("Inter-Bold", inter.BoldFont)
		if err != nil {
			return err
		}

		writeLogo(&pdf, file.Logo, file.From)
		writeTitle(&pdf, file.Title, file.Id, file.Date)
		writeBillTo(&pdf, file.To)
		writeHeaderRow(&pdf)
		subtotal := 0.0

		for i := range file.Items {
			quantity := 1.0
			if len(file.Quantities) > i {
				quantity = file.Quantities[i]
			}

			rate := 0.0
			if len(file.Rates) > i {
				rate = file.Rates[i]
			}

			writeRow(&pdf, file.Items[i], quantity, rate, file.Currency)
			subtotal += float64(quantity) * rate
		}

		if file.Note != "" {
			writeNotes(&pdf, file.Note)
		}

		writeTotals(
			&pdf,
			subtotal,
			subtotal*file.Tax,
			subtotal*file.Discount,
			file.ForeignExchange.Rate,
			file.Currency,
			file.ForeignExchange.Currency,
		)

		if file.Due != "" {
			writeDueDate(&pdf, file.Due)
		}

		writeFooter(&pdf, file.Id)
		output := strings.TrimSuffix(file.Output, ".pdf") + ".pdf"
		err = pdf.WritePdf(output)
		if err != nil {
			return err
		}

		fmt.Printf("Generated %s\n", output)

		return nil
	}
}

func Command() *cobra.Command {
	file := invoice{}
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate an invoice",
		Long:  `Generate an invoice`,
		RunE:  generate(&file),
	}

	defaultInvoice := getDefaultInvoice()

	generateCmd.Flags().StringVar(&file.Id, "id", time.Now().Format("20060102"), "ID")
	generateCmd.Flags().StringVar(&file.Title, "title", "INVOICE", "Title")

	generateCmd.Flags().Float64SliceVarP(&file.Rates, "rate", "r", defaultInvoice.Rates, "Rates")
	generateCmd.Flags().
		Float64SliceVarP(&file.Quantities, "quantity", "q", defaultInvoice.Quantities, "Quantities")
	generateCmd.Flags().StringSliceVarP(&file.Items, "item", "i", defaultInvoice.Items, "Items")

	generateCmd.Flags().StringVarP(&file.Logo, "logo", "l", defaultInvoice.Logo, "Company logo")
	generateCmd.Flags().StringVarP(&file.From, "from", "f", defaultInvoice.From, "Issuing company")
	generateCmd.Flags().StringVarP(&file.To, "to", "t", defaultInvoice.To, "Recipient company")
	generateCmd.Flags().StringVar(&file.Date, "date", defaultInvoice.Date, "Date")
	generateCmd.Flags().StringVar(&file.Due, "due", defaultInvoice.Due, "Payment due date")

	generateCmd.Flags().Float64Var(&file.Tax, "tax", defaultInvoice.Tax, "Tax (Percentage)")
	generateCmd.Flags().
		Float64VarP(&file.Discount, "discount", "d", defaultInvoice.Discount, "Discount (Percentage)")
	generateCmd.Flags().
		StringVarP(&file.Currency, "currency", "c", defaultInvoice.Currency, "Currency")

	generateCmd.Flags().
		Float64Var(&file.ForeignExchange.Rate, "exchange-rate", defaultInvoice.ForeignExchange.Rate, "Foreign Exchange Rate")
	generateCmd.Flags().
		StringVar(&file.ForeignExchange.Currency, "exchange-currency", defaultInvoice.ForeignExchange.Currency, "Foreign Exchange Currency")

	generateCmd.Flags().StringVarP(&file.Note, "note", "n", "", "Note")
	generateCmd.Flags().StringVarP(&file.Output, "output", "o", "invoice.pdf", "Output file (.pdf)")

	return generateCmd
}
