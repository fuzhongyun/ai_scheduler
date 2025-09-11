package constants

type KnowledgeId string

const (
	KnowledgeIdZltx KnowledgeId = "kb-00000001"
)

func (k KnowledgeId) String() string {
	return string(k)
}

var CallerKnowledgeIdMap = map[Caller]KnowledgeId{
	CallerZltx: KnowledgeIdZltx,
}

func GetKnowledgeId(caller Caller) KnowledgeId {
	if _, ok := CallerKnowledgeIdMap[caller]; !ok {
		return ""
	}
	return CallerKnowledgeIdMap[caller]
}
