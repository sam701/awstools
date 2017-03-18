package cf

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type resourceTag struct {
	key, value string
}

func tagFromString(str string) *resourceTag {
	parts := strings.Split(str, ":")
	return &resourceTag{parts[0], strings.ToLower(parts[1])}
}

type tagList []*resourceTag

func tagListFromStrings(ar []string) tagList {
	res := []*resourceTag{}
	for _, tagDef := range ar {
		res = append(res, tagFromString(tagDef))
	}
	return res
}

func (s tagList) match(cfTags []*cloudformation.Tag) bool {
	for _, tag := range s {
		foundMatch := false
		for _, cfTag := range cfTags {
			if tag.key == *cfTag.Key && strings.Contains(strings.ToLower(*cfTag.Value), tag.value) {
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			return false
		}
	}
	return true
}

type stackFilter struct {
	name        string
	namePattern string
	tags        tagList
}

func createStackFilter(name, namePattern string, tags []string) *stackFilter {
	return &stackFilter{name, strings.ToLower(namePattern), tagListFromStrings(tags)}
}

func (s *stackFilter) match(stack *cloudformation.Stack) bool {
	if s.name != "" {
		return s.name == *stack.StackName
	}

	if s.namePattern != "" {
		if !strings.Contains(*stack.StackName, s.namePattern) {
			return false
		}
	}

	if len(s.tags) > 0 {
		if !s.tags.match(stack.Tags) {
			return false
		}
	}

	return true
}
