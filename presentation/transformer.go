package presentation

type (
	Transformable interface {
		Transform() interface{}
	}
)

func TransformList[T Transformable](list []T) []interface{} {
	result := []interface{}{}
	for _, v := range list {
		result = append(result, v.Transform())
	}
	return result
}

func TransformListAny(list []interface{}) []interface{} {
	result := []interface{}{}
	for _, v := range list {
		if transformer, ok := v.(Transformable); ok {
			result = append(result, transformer.Transform())
		} else {
			result = append(result, v)
		}
	}
	return result
}
