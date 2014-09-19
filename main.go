
package rule

type Error map[string][]string
type T func(obj interface{}) Error
type Msgs []string

func Passed(err Error) bool {
	return len(err) == 0
}

func Failed(err Error) bool {
	return !Passed(err)
}

func One(rules ...T) T {
	return T(func(obj interface{}) Error {
		for _,rule := range rules {
			err := rule(obj)
			if Failed(err) {
				return err
			}
		}
		return nil
	})
}

func All(rules ...T) T {
	return T(func(obj interface{}) Error {
		allErrors := make(Error)
		for _,rule := range rules {
			err := rule(obj)
			Merge(allErrors, err)
		}
		return allErrors
	})
}

func Merge(dest Error, src Error) Error {
	if dest == nil {
		return src
	}
	for k,v := range src {
		dest[k] = append(dest[k], v...)
	}
	return dest
}

func AnError(name string, messages ...string) Error {
	return Error{name : Msgs(messages)}
}

func main() {
}



