package splitbatchprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/pdata"
)

func TestSplitDifferentTracesIntoDifferentBatches(t *testing.T) {
	// we have 1 ResourceSpans with 1 ILS and two traceIDs, resulting in two batches
	rs := pdata.NewResourceSpans()
	rs.InitEmpty()
	rs.InstrumentationLibrarySpans().Resize(1)

	// the first ILS has two spans
	ils := rs.InstrumentationLibrarySpans().At(0)
	library := ils.InstrumentationLibrary()
	library.InitEmpty()
	library.SetName("first-library")
	ils.Spans().Resize(2)
	firstSpan := ils.Spans().At(0)
	firstSpan.SetName("first-batch-first-span")
	firstSpan.SetTraceID(pdata.NewTraceID([]byte{1, 2, 3, 4}))
	secondSpan := ils.Spans().At(1)
	secondSpan.SetName("first-batch-second-span")
	secondSpan.SetTraceID(pdata.NewTraceID([]byte{2, 3, 4, 5}))

	inBatch := pdata.NewTraces()
	inBatch.ResourceSpans().Append(rs)

	// test
	next := &componenttest.ExampleExporterConsumer{}
	processor := newSplitBatch(next)
	err := processor.ConsumeTraces(context.Background(), inBatch)

	// verify
	assert.NoError(t, err)
	assert.Len(t, next.Traces, 2)

	// first batch
	firstOutILS := next.Traces[0].ResourceSpans().At(0).InstrumentationLibrarySpans().At(0)
	assert.Equal(t, library.Name(), firstOutILS.InstrumentationLibrary().Name())
	assert.Equal(t, firstSpan.Name(), firstOutILS.Spans().At(0).Name())

	// second batch
	secondOutILS := next.Traces[1].ResourceSpans().At(0).InstrumentationLibrarySpans().At(0)
	assert.Equal(t, library.Name(), secondOutILS.InstrumentationLibrary().Name())
	assert.Equal(t, secondSpan.Name(), secondOutILS.Spans().At(0).Name())
}
