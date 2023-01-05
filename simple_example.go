package mapstruct

// Check field name.
// Field description by string.
// 1. Генерируем состав полей
// 2. Заполняем пользовательские настройки
// 3. Заполняем настройки по умолчанию
// 4. Интерфейсы не преобразовываем

// go:generate mapstruct --srcPkg path --srcType TypeName dstPkg  path --dstType TypeName
// --options srcOption
// 1. Рекурсивная загрузка пакетов
// 2. Структура катологов при повторении имен типов в пакете

// Парсинг пакетов -> Парсинг настроек -> Формирование промежуточного представления (FieldDefinition) -> генерация по макету

type Source struct {
	SimpleField string
	Interface   error

	ComplexFieldTransform innerStructSrc
	ComplexField          innerStructSrc
}

type innerStructSrc struct {
	SimpleField int
	Array       []string
}

type Distination struct {
	SimpleField  string
	Interface    error
	Array        []string
	ComplexField innerStructDst
}

type innerStructDst struct {
	SimpleField string
	Array       []string
}

type Converter interface {
	Map(*Source) Distination
}

type Settings struct{}

type MapStrOption func(s *Settings)

//var mapSettings = func() []MapStrOption {
//	return []MapStrOption{
//		Field("ComplexFieldTransform").Map(
//			Field("Array").To("Array")
//		)
//	}
//}
