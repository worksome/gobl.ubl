package ubl_test

import (
	"testing"

	ubl "github.com/invopop/gobl.ubl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseInvoiceTypes(t *testing.T) {
	t.Run("standard invoice (380)", func(t *testing.T) {
		e, err := testParseInvoice("peppol/base-example.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.Equal(t, bill.InvoiceTypeStandard, inv.Type)
		assert.Empty(t, inv.Tags)
	})

	t.Run("credit note (381)", func(t *testing.T) {
		e, err := testParseInvoice("peppol/base-creditnote-correction.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.Equal(t, bill.InvoiceTypeCreditNote, inv.Type)
		assert.Empty(t, inv.Tags)
	})

	t.Run("proforma invoice (325)", func(t *testing.T) {
		e, err := testParseInvoice("peppol/proforma-invoice.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.Equal(t, bill.InvoiceTypeProforma, inv.Type)
		assert.Empty(t, inv.Tags)
	})

	t.Run("self-billed invoice (389)", func(t *testing.T) {
		e, err := testParseInvoice("peppol/self-billed-invoice.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.Equal(t, bill.InvoiceTypeStandard, inv.Type)
		assert.True(t, inv.HasTags(tax.TagSelfBilled), "should have self-billed tag")
	})

	t.Run("partial invoice (326)", func(t *testing.T) {
		e, err := testParseInvoice("peppol/partial-invoice.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.Equal(t, bill.InvoiceTypeStandard, inv.Type)
		assert.True(t, inv.HasTags(tax.TagPartial), "should have partial tag")
	})

	t.Run("self-billed credit note (261)", func(t *testing.T) {
		e, err := testParseInvoice("peppol/self-billed-creditnote.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.Equal(t, bill.InvoiceTypeCreditNote, inv.Type)
		assert.True(t, inv.HasTags(tax.TagSelfBilled), "should have self-billed tag")
	})
}

func TestParseInvoiceTags(t *testing.T) {
	t.Run("invoice with self-billed tag", func(t *testing.T) {
		e, err := testParseInvoice("peppol/self-billed-invoice.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.True(t, inv.HasTags(tax.TagSelfBilled), "should have self-billed tag")
	})

	t.Run("invoice with partial tag", func(t *testing.T) {
		e, err := testParseInvoice("peppol/partial-invoice.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.True(t, inv.HasTags(tax.TagPartial), "should have partial tag")
	})

	t.Run("credit note with self-billed tag", func(t *testing.T) {
		e, err := testParseInvoice("peppol/self-billed-creditnote.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.True(t, inv.HasTags(tax.TagSelfBilled), "should have self-billed tag")
	})

	t.Run("standard invoice without tags", func(t *testing.T) {
		e, err := testParseInvoice("peppol/base-example.xml")
		require.NoError(t, err)

		inv, ok := e.Extract().(*bill.Invoice)
		require.True(t, ok)

		assert.False(t, inv.HasTags(tax.TagSelfBilled), "standard invoice should not have self-billed tag")
		assert.False(t, inv.HasTags(tax.TagPartial), "standard invoice should not have partial tag")
	})
}

func TestParseInvoiceTypeAndTagCombinations(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		expectedType string
		expectedTags []string
	}{
		{
			name:         "standard invoice (380)",
			filename:     "peppol/base-example.xml",
			expectedType: string(bill.InvoiceTypeStandard),
			expectedTags: nil,
		},
		{
			name:         "credit note (381)",
			filename:     "peppol/base-creditnote-correction.xml",
			expectedType: string(bill.InvoiceTypeCreditNote),
			expectedTags: nil,
		},
		{
			name:         "proforma (325)",
			filename:     "peppol/proforma-invoice.xml",
			expectedType: string(bill.InvoiceTypeProforma),
			expectedTags: nil,
		},
		{
			name:         "self-billed standard (389)",
			filename:     "peppol/self-billed-invoice.xml",
			expectedType: string(bill.InvoiceTypeStandard),
			expectedTags: []string{string(tax.TagSelfBilled)},
		},
		{
			name:         "partial standard (326)",
			filename:     "peppol/partial-invoice.xml",
			expectedType: string(bill.InvoiceTypeStandard),
			expectedTags: []string{string(tax.TagPartial)},
		},
		{
			name:         "self-billed credit note (261)",
			filename:     "peppol/self-billed-creditnote.xml",
			expectedType: string(bill.InvoiceTypeCreditNote),
			expectedTags: []string{string(tax.TagSelfBilled)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := testParseInvoice(tt.filename)
			require.NoError(t, err)

			inv, ok := e.Extract().(*bill.Invoice)
			require.True(t, ok)

			assert.Equal(t, tt.expectedType, string(inv.Type), "invoice type mismatch")

			if tt.expectedTags == nil {
				assert.False(t, inv.HasTags(tax.TagSelfBilled), "should not have self-billed tag")
				assert.False(t, inv.HasTags(tax.TagPartial), "should not have partial tag")
			} else {
				for _, tag := range tt.expectedTags {
					assert.True(t, inv.HasTags(cbc.Key(tag)), "missing expected tag: %s", tag)
				}
			}
		})
	}
}

func TestParseUUID(t *testing.T) {
	t.Run("parses UUID from OIOUBL XML", func(t *testing.T) {
		xmlInput := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Invoice xmlns:cac="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2" xmlns:cbc="urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2" xmlns="urn:oasis:names:specification:ubl:schema:xsd:Invoice-2">
  <cbc:CustomizationID>urn:fdc:oioubl.dk:trns:billing:invoice:3.0</cbc:CustomizationID>
  <cbc:ProfileID>urn:fdc:oioubl.dk:bis:billing_with_response:3</cbc:ProfileID>
  <cbc:ID>TEST-UUID-001</cbc:ID>
  <cbc:UUID>019cde16-3215-75a1-84f9-4d410281281f</cbc:UUID>
  <cbc:IssueDate>2026-01-01</cbc:IssueDate>
  <cbc:InvoiceTypeCode>380</cbc:InvoiceTypeCode>
  <cbc:DocumentCurrencyCode>DKK</cbc:DocumentCurrencyCode>
  <cac:AccountingSupplierParty>
    <cac:Party>
      <cac:PartyLegalEntity>
        <cbc:RegistrationName>Test Supplier</cbc:RegistrationName>
      </cac:PartyLegalEntity>
    </cac:Party>
  </cac:AccountingSupplierParty>
  <cac:AccountingCustomerParty>
    <cac:Party>
      <cac:PartyLegalEntity>
        <cbc:RegistrationName>Test Customer</cbc:RegistrationName>
      </cac:PartyLegalEntity>
    </cac:Party>
  </cac:AccountingCustomerParty>
  <cac:LegalMonetaryTotal>
    <cbc:PayableAmount currencyID="DKK">100.00</cbc:PayableAmount>
  </cac:LegalMonetaryTotal>
</Invoice>`)

		parsed, err := ubl.Parse(xmlInput)
		require.NoError(t, err)

		inv, ok := parsed.(*ubl.Invoice)
		require.True(t, ok)
		assert.Equal(t, "019cde16-3215-75a1-84f9-4d410281281f", inv.UUID)

		env, err := inv.Convert()
		require.NoError(t, err)

		goblInv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Equal(t, "019cde16-3215-75a1-84f9-4d410281281f", goblInv.UUID.String(),
			"UUID from UBL XML should be preserved in GOBL invoice")
	})

	t.Run("round-trips UUID through OIOUBL30", func(t *testing.T) {
		doc, err := testInvoiceFrom("oioubl30-invoice-example.json")
		require.NoError(t, err)
		require.NotEmpty(t, doc.UUID)

		data, err := ubl.Bytes(doc)
		require.NoError(t, err)

		parsed, err := ubl.Parse(data)
		require.NoError(t, err)

		inv, ok := parsed.(*ubl.Invoice)
		require.True(t, ok)
		assert.Equal(t, doc.UUID, inv.UUID, "UUID should survive XML round-trip")

		env, err := inv.Convert()
		require.NoError(t, err)

		goblInv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Equal(t, doc.UUID, goblInv.UUID.String(),
			"UUID should survive full round-trip through GOBL")
	})
}

func TestParseCopyIndicator(t *testing.T) {
	t.Run("true CopyIndicator is stored in meta", func(t *testing.T) {
		xmlInput := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Invoice xmlns:cac="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2" xmlns:cbc="urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2" xmlns="urn:oasis:names:specification:ubl:schema:xsd:Invoice-2">
  <cbc:ID>TEST-COPY</cbc:ID>
  <cbc:CopyIndicator>true</cbc:CopyIndicator>
  <cbc:IssueDate>2026-01-01</cbc:IssueDate>
  <cbc:InvoiceTypeCode>380</cbc:InvoiceTypeCode>
  <cbc:DocumentCurrencyCode>DKK</cbc:DocumentCurrencyCode>
  <cac:AccountingSupplierParty>
    <cac:Party>
      <cac:PartyLegalEntity>
        <cbc:RegistrationName>Test</cbc:RegistrationName>
      </cac:PartyLegalEntity>
    </cac:Party>
  </cac:AccountingSupplierParty>
  <cac:AccountingCustomerParty>
    <cac:Party>
      <cac:PartyLegalEntity>
        <cbc:RegistrationName>Test</cbc:RegistrationName>
      </cac:PartyLegalEntity>
    </cac:Party>
  </cac:AccountingCustomerParty>
  <cac:LegalMonetaryTotal>
    <cbc:PayableAmount currencyID="DKK">100.00</cbc:PayableAmount>
  </cac:LegalMonetaryTotal>
</Invoice>`)

		parsed, err := ubl.Parse(xmlInput)
		require.NoError(t, err)

		inv, ok := parsed.(*ubl.Invoice)
		require.True(t, ok)
		assert.True(t, inv.CopyIndicator)

		env, err := inv.Convert()
		require.NoError(t, err)

		goblInv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Equal(t, "true", goblInv.Meta["copy"],
			"CopyIndicator should be preserved as meta key 'copy'")
	})

	t.Run("false CopyIndicator does not set meta", func(t *testing.T) {
		xmlInput := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Invoice xmlns:cac="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2" xmlns:cbc="urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2" xmlns="urn:oasis:names:specification:ubl:schema:xsd:Invoice-2">
  <cbc:ID>TEST-NO-COPY</cbc:ID>
  <cbc:CopyIndicator>false</cbc:CopyIndicator>
  <cbc:IssueDate>2026-01-01</cbc:IssueDate>
  <cbc:InvoiceTypeCode>380</cbc:InvoiceTypeCode>
  <cbc:DocumentCurrencyCode>DKK</cbc:DocumentCurrencyCode>
  <cac:AccountingSupplierParty>
    <cac:Party>
      <cac:PartyLegalEntity>
        <cbc:RegistrationName>Test</cbc:RegistrationName>
      </cac:PartyLegalEntity>
    </cac:Party>
  </cac:AccountingSupplierParty>
  <cac:AccountingCustomerParty>
    <cac:Party>
      <cac:PartyLegalEntity>
        <cbc:RegistrationName>Test</cbc:RegistrationName>
      </cac:PartyLegalEntity>
    </cac:Party>
  </cac:AccountingCustomerParty>
  <cac:LegalMonetaryTotal>
    <cbc:PayableAmount currencyID="DKK">100.00</cbc:PayableAmount>
  </cac:LegalMonetaryTotal>
</Invoice>`)

		parsed, err := ubl.Parse(xmlInput)
		require.NoError(t, err)

		inv, ok := parsed.(*ubl.Invoice)
		require.True(t, ok)

		env, err := inv.Convert()
		require.NoError(t, err)

		goblInv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		_, hasCopy := goblInv.Meta["copy"]
		assert.False(t, hasCopy, "false CopyIndicator should not create meta entry")
	})
}
