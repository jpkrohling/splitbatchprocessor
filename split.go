package splitbatchprocessor

import "go.opentelemetry.io/collector/consumer/pdata"

func splitByTrace(rs pdata.ResourceSpans) []pdata.ResourceSpans {
	// for each span in the resource spans, we group them into batches of rs/ils/traceID.
	// if the same traceID exists in different ils, they land in different batches.
	var result []pdata.ResourceSpans

	for i := 0; i < rs.InstrumentationLibrarySpans().Len(); i++ {
		// the batches for this ILS
		batches := map[string]pdata.ResourceSpans{}

		ils := rs.InstrumentationLibrarySpans().At(i)
		for j := 0; j < ils.Spans().Len(); j++ {
			span := ils.Spans().At(j)
			key := string(span.TraceID().Bytes())

			// for the first traceID in the ILS, initialize the map entry
			// and add the singleTraceBatch to the result list
			if _, ok := batches[key]; !ok {
				newRS := pdata.NewResourceSpans()
				newRS.InitEmpty()
				// currently, the ResourceSpans implementation has only a Resource and an ILS. We'll copy the Resource
				// and set our own ILS
				rs.Resource().CopyTo(newRS.Resource())

				newILS := pdata.NewInstrumentationLibrarySpans()
				newILS.InitEmpty()
				// currently, the ILS implementation has only an InstrumentationLibrary and spans. We'll copy the library
				// and set our own spans
				ils.InstrumentationLibrary().CopyTo(newILS.InstrumentationLibrary())
				newRS.InstrumentationLibrarySpans().Append(newILS)
				batches[key] = newRS

				result = append(result, newRS)
			}

			// there is only one instrumentation library per batch
			batches[key].InstrumentationLibrarySpans().At(0).Spans().Append(span)
		}
	}

	return result
}
