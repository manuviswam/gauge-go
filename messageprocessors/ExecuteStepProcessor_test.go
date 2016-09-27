package messageprocessors

import (
	"testing"

	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"github.com/manuviswam/gauge-go/models"
	"github.com/stretchr/testify/assert"
)

func TestShouldRunStep(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	called := false
	context := &t.GaugeContext{
		Steps: []t.Step{t.Step{
			Description: stepText,
			Impl:        func(args ...interface{}) { called = true },
		},
		},
	}

	msg := &m.Message{
		MessageType: m.Message_ExecuteStep.Enum(),
		MessageId:   &msgId,
		ExecuteStepRequest: &m.ExecuteStepRequest{
			ParsedStepText: &stepText,
		},
	}

	p := ExecuteStepProcessor{}

	p.Process(msg, context)

	assert.True(tst, called)

}

func TestShouldRunReturnExecutionStatusResponseWithSameId(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	called := false
	context := &t.GaugeContext{
		Steps: []t.Step{t.Step{
			Description: stepText,
			Impl:        func(args ...interface{}) { called = true },
		},
		},
	}

	msg := &m.Message{
		MessageType: m.Message_ExecuteStep.Enum(),
		MessageId:   &msgId,
		ExecuteStepRequest: &m.ExecuteStepRequest{
			ParsedStepText: &stepText,
		},
	}

	p := ExecuteStepProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}


func TestShouldRunStepWithTableParam(tst *testing.T) {
	stepText := "Step description <table>"
	msgId := int64(12345)
	called := false
	context := &t.GaugeContext{
		Steps: []t.Step{t.Step{
			Description: stepText,
			Impl:     func(tbl *models.Table) { called = true },
		},
		},
	}

	headers := []string{"Header 1", "Header 2"}
	columns := []string{"Value 1", "Value 2"}
	row := &m.ProtoTableRow{Cells: columns}
	rows := []*m.ProtoTableRow{row}
	p := &m.ProtoTable{
		Headers: &m.ProtoTableRow{
			Cells: headers,
		},
		Rows: rows,
	}

	msg := &m.Message{
		MessageType: m.Message_ExecuteStep.Enum(),
		MessageId:   &msgId,
		ExecuteStepRequest: &m.ExecuteStepRequest{
			ParsedStepText: &stepText,
			Parameters: []*m.Parameter{
				&m.Parameter{
					ParameterType: m.Parameter_Table.Enum(),
					Table: p,
				},
			},
		},
	}

	proc := ExecuteStepProcessor{}
	proc.Process(msg, context)

	assert.True(tst, called)
}

func TestShouldRunStepWithSpecialTableParam(tst *testing.T) {
	stepText := "Step description <table>"
	msgId := int64(12345)
	called := false
	headers := []string{"Header 1", "Header 2"}
	columns := []string{"Value 1", "Value 2"}
	context := &t.GaugeContext{
		Steps: []t.Step{t.Step{
			Description: stepText,
			Impl:     func(tbl *models.Table) {
				called = true
				assert.Equal(tst, tbl.Rows[0].Cells[0], columns[0])
				assert.Equal(tst, tbl.Rows[0].Cells[1], columns[1])
			 },
		},
		},
	}

	row := &m.ProtoTableRow{Cells: columns}
	rows := []*m.ProtoTableRow{row}

	p := &m.ProtoTable{
		Headers: &m.ProtoTableRow{
			Cells: headers,
		},
		Rows: rows,
	}

	msg := &m.Message{
		MessageType: m.Message_ExecuteStep.Enum(),
		MessageId:   &msgId,
		ExecuteStepRequest: &m.ExecuteStepRequest{
			ParsedStepText: &stepText,
			Parameters: []*m.Parameter{
				&m.Parameter{
					ParameterType: m.Parameter_Special_Table.Enum(),
					Table: p,
				},
			},
		},
	}

	proc := ExecuteStepProcessor{}
	proc.Process(msg, context)

	assert.True(tst, called)
}
