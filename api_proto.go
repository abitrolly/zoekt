// Copyright 2016 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zoekt // import "github.com/sourcegraph/zoekt"

import (
	"math/rand"
	"reflect"

	v1 "github.com/sourcegraph/zoekt/grpc/v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FileMatchFromProto(p *v1.FileMatch) FileMatch {
	lineMatches := make([]LineMatch, len(p.GetLineMatches()))
	for i, lineMatch := range p.GetLineMatches() {
		lineMatches[i] = LineMatchFromProto(lineMatch)
	}

	chunkMatches := make([]ChunkMatch, len(p.GetChunkMatches()))
	for i, chunkMatch := range p.GetChunkMatches() {
		chunkMatches[i] = ChunkMatchFromProto(chunkMatch)
	}

	return FileMatch{
		Score:              p.GetScore(),
		Debug:              p.GetDebug(),
		FileName:           p.GetFileName(),
		Repository:         p.GetRepository(),
		Branches:           p.GetBranches(),
		LineMatches:        lineMatches,
		ChunkMatches:       chunkMatches,
		RepositoryID:       p.GetRepositoryId(),
		RepositoryPriority: p.GetRepositoryPriority(),
		Content:            p.GetContent(),
		Checksum:           p.GetChecksum(),
		Language:           p.GetLanguage(),
		SubRepositoryName:  p.GetSubRepositoryName(),
		SubRepositoryPath:  p.GetSubRepositoryPath(),
		Version:            p.GetVersion(),
	}
}

func (m *FileMatch) ToProto() *v1.FileMatch {
	lineMatches := make([]*v1.LineMatch, len(m.LineMatches))
	for i, lm := range m.LineMatches {
		lineMatches[i] = lm.ToProto()
	}

	chunkMatches := make([]*v1.ChunkMatch, len(m.ChunkMatches))
	for i, cm := range m.ChunkMatches {
		chunkMatches[i] = cm.ToProto()
	}

	return &v1.FileMatch{
		Score:              m.Score,
		Debug:              m.Debug,
		FileName:           m.FileName,
		Repository:         m.Repository,
		Branches:           m.Branches,
		LineMatches:        lineMatches,
		ChunkMatches:       chunkMatches,
		RepositoryId:       m.RepositoryID,
		RepositoryPriority: m.RepositoryPriority,
		Content:            m.Content,
		Checksum:           m.Checksum,
		Language:           m.Language,
		SubRepositoryName:  m.SubRepositoryName,
		SubRepositoryPath:  m.SubRepositoryPath,
		Version:            m.Version,
	}
}

func ChunkMatchFromProto(p *v1.ChunkMatch) ChunkMatch {
	ranges := make([]Range, len(p.GetRanges()))
	for i, r := range p.GetRanges() {
		ranges[i] = RangeFromProto(r)
	}

	symbols := make([]*Symbol, len(p.GetSymbolInfo()))
	for i, r := range p.GetSymbolInfo() {
		symbols[i] = SymbolFromProto(r)
	}

	return ChunkMatch{
		Content:      p.GetContent(),
		ContentStart: LocationFromProto(p.GetContentStart()),
		FileName:     p.GetFileName(),
		Ranges:       ranges,
		SymbolInfo:   symbols,
		Score:        p.GetScore(),
		DebugScore:   p.GetDebugScore(),
	}
}

func (cm *ChunkMatch) ToProto() *v1.ChunkMatch {
	ranges := make([]*v1.Range, len(cm.Ranges))
	for i, r := range cm.Ranges {
		ranges[i] = r.ToProto()
	}

	symbolInfo := make([]*v1.SymbolInfo, len(cm.SymbolInfo))
	for i, si := range cm.SymbolInfo {
		symbolInfo[i] = si.ToProto()
	}

	return &v1.ChunkMatch{
		Content:      cm.Content,
		ContentStart: cm.ContentStart.ToProto(),
		FileName:     cm.FileName,
		Ranges:       ranges,
		SymbolInfo:   symbolInfo,
		Score:        cm.Score,
		DebugScore:   cm.DebugScore,
	}
}

func RangeFromProto(p *v1.Range) Range {
	return Range{
		Start: LocationFromProto(p.GetStart()),
		End:   LocationFromProto(p.GetEnd()),
	}
}

func (r *Range) ToProto() *v1.Range {
	return &v1.Range{
		Start: r.Start.ToProto(),
		End:   r.End.ToProto(),
	}
}

func LocationFromProto(p *v1.Location) Location {
	return Location{
		ByteOffset: p.GetByteOffset(),
		LineNumber: p.GetLineNumber(),
		Column:     p.GetColumn(),
	}
}

func (l *Location) ToProto() *v1.Location {
	return &v1.Location{
		ByteOffset: l.ByteOffset,
		LineNumber: l.LineNumber,
		Column:     l.Column,
	}
}

func LineMatchFromProto(p *v1.LineMatch) LineMatch {
	lineFragments := make([]LineFragmentMatch, len(p.GetLineFragments()))
	for i, lineFragment := range p.GetLineFragments() {
		lineFragments[i] = LineFragmentMatchFromProto(lineFragment)
	}

	return LineMatch{
		Line:          p.GetLine(),
		LineStart:     int(p.GetLineStart()),
		LineEnd:       int(p.GetLineEnd()),
		LineNumber:    int(p.GetLineNumber()),
		Before:        p.GetBefore(),
		After:         p.GetAfter(),
		FileName:      p.GetFileName(),
		Score:         p.GetScore(),
		DebugScore:    p.GetDebugScore(),
		LineFragments: lineFragments,
	}
}

func (lm *LineMatch) ToProto() *v1.LineMatch {
	fragments := make([]*v1.LineFragmentMatch, len(lm.LineFragments))
	for i, fragment := range lm.LineFragments {
		fragments[i] = fragment.ToProto()
	}

	return &v1.LineMatch{
		Line:          lm.Line,
		LineStart:     int64(lm.LineStart),
		LineEnd:       int64(lm.LineEnd),
		LineNumber:    int64(lm.LineNumber),
		Before:        lm.Before,
		After:         lm.After,
		FileName:      lm.FileName,
		Score:         lm.Score,
		DebugScore:    lm.DebugScore,
		LineFragments: fragments,
	}
}

func SymbolFromProto(p *v1.SymbolInfo) *Symbol {
	if p == nil {
		return nil
	}

	return &Symbol{
		Sym:        p.GetSym(),
		Kind:       p.GetKind(),
		Parent:     p.GetParent(),
		ParentKind: p.GetParentKind(),
	}
}

func (s *Symbol) ToProto() *v1.SymbolInfo {
	if s == nil {
		return nil
	}

	return &v1.SymbolInfo{
		Sym:        s.Sym,
		Kind:       s.Kind,
		Parent:     s.Parent,
		ParentKind: s.ParentKind,
	}
}

func LineFragmentMatchFromProto(p *v1.LineFragmentMatch) LineFragmentMatch {
	return LineFragmentMatch{
		LineOffset:  int(p.GetLineOffset()),
		Offset:      p.GetOffset(),
		MatchLength: int(p.GetMatchLength()),
		SymbolInfo:  SymbolFromProto(p.GetSymbolInfo()),
	}
}

func (lfm *LineFragmentMatch) ToProto() *v1.LineFragmentMatch {
	return &v1.LineFragmentMatch{
		LineOffset:  int64(lfm.LineOffset),
		Offset:      lfm.Offset,
		MatchLength: int64(lfm.MatchLength),
		SymbolInfo:  lfm.SymbolInfo.ToProto(),
	}
}

func FlushReasonFromProto(p v1.FlushReason) FlushReason {
	switch p {
	case v1.FlushReason_TIMER_EXPIRED:
		return FlushReasonTimerExpired
	case v1.FlushReason_FINAL_FLUSH:
		return FlushReasonFinalFlush
	case v1.FlushReason_MAX_SIZE:
		return FlushReasonMaxSize
	default:
		return FlushReason(0)
	}
}

func (fr FlushReason) ToProto() v1.FlushReason {
	switch fr {
	case FlushReasonTimerExpired:
		return v1.FlushReason_TIMER_EXPIRED
	case FlushReasonFinalFlush:
		return v1.FlushReason_FINAL_FLUSH
	case FlushReasonMaxSize:
		return v1.FlushReason_MAX_SIZE
	default:
		return v1.FlushReason_UNKNOWN
	}
}

// Generate valid reasons for quickchecks
func (fr FlushReason) Generate(rand *rand.Rand, size int) reflect.Value {
	switch rand.Int() % 4 {
	case 1:
		return reflect.ValueOf(FlushReasonMaxSize)
	case 2:
		return reflect.ValueOf(FlushReasonFinalFlush)
	case 3:
		return reflect.ValueOf(FlushReasonTimerExpired)
	default:
		return reflect.ValueOf(FlushReason(0))
	}
}

func StatsFromProto(p *v1.Stats) Stats {
	return Stats{
		ContentBytesLoaded:   p.GetContentBytesLoaded(),
		IndexBytesLoaded:     p.GetIndexBytesLoaded(),
		Crashes:              int(p.GetCrashes()),
		Duration:             p.GetDuration().AsDuration(),
		FileCount:            int(p.GetFileCount()),
		ShardFilesConsidered: int(p.GetShardFilesConsidered()),
		FilesConsidered:      int(p.GetFilesConsidered()),
		FilesLoaded:          int(p.GetFilesLoaded()),
		FilesSkipped:         int(p.GetFilesSkipped()),
		ShardsScanned:        int(p.GetShardsScanned()),
		ShardsSkipped:        int(p.GetShardsSkipped()),
		ShardsSkippedFilter:  int(p.GetShardsSkippedFilter()),
		MatchCount:           int(p.GetMatchCount()),
		NgramMatches:         int(p.GetNgramMatches()),
		Wait:                 p.GetWait().AsDuration(),
		RegexpsConsidered:    int(p.GetRegexpsConsidered()),
		FlushReason:          FlushReasonFromProto(p.GetFlushReason()),
	}
}

func (s *Stats) ToProto() *v1.Stats {
	return &v1.Stats{
		ContentBytesLoaded:   s.ContentBytesLoaded,
		IndexBytesLoaded:     s.IndexBytesLoaded,
		Crashes:              int64(s.Crashes),
		Duration:             durationpb.New(s.Duration),
		FileCount:            int64(s.FileCount),
		ShardFilesConsidered: int64(s.ShardFilesConsidered),
		FilesConsidered:      int64(s.FilesConsidered),
		FilesLoaded:          int64(s.FilesLoaded),
		FilesSkipped:         int64(s.FilesSkipped),
		ShardsScanned:        int64(s.ShardsScanned),
		ShardsSkipped:        int64(s.ShardsSkipped),
		ShardsSkippedFilter:  int64(s.ShardsSkippedFilter),
		MatchCount:           int64(s.MatchCount),
		NgramMatches:         int64(s.NgramMatches),
		Wait:                 durationpb.New(s.Wait),
		RegexpsConsidered:    int64(s.RegexpsConsidered),
		FlushReason:          s.FlushReason.ToProto(),
	}
}

func ProgressFromProto(p *v1.Progress) Progress {
	return Progress{
		Priority:           p.GetPriority(),
		MaxPendingPriority: p.GetMaxPendingPriority(),
	}
}

func (p *Progress) ToProto() *v1.Progress {
	return &v1.Progress{
		Priority:           p.Priority,
		MaxPendingPriority: p.MaxPendingPriority,
	}
}

func SearchResultFromProto(p *v1.SearchResponse) *SearchResult {
	if p == nil {
		return nil
	}

	files := make([]FileMatch, len(p.GetFiles()))
	for i, file := range p.GetFiles() {
		files[i] = FileMatchFromProto(file)
	}

	return &SearchResult{
		Stats:         StatsFromProto(p.GetStats()),
		Progress:      ProgressFromProto(p.GetProgress()),
		Files:         files,
		RepoURLs:      p.RepoUrls,
		LineFragments: p.LineFragments,
	}
}

func (sr *SearchResult) ToProto() *v1.SearchResponse {
	if sr == nil {
		return nil
	}

	files := make([]*v1.FileMatch, len(sr.Files))
	for i, file := range sr.Files {
		files[i] = file.ToProto()
	}

	return &v1.SearchResponse{
		Stats:         sr.Stats.ToProto(),
		Progress:      sr.Progress.ToProto(),
		Files:         files,
		RepoUrls:      sr.RepoURLs,
		LineFragments: sr.LineFragments,
	}
}

func RepositoryBranchFromProto(p *v1.RepositoryBranch) RepositoryBranch {
	return RepositoryBranch{
		Name:    p.GetName(),
		Version: p.GetVersion(),
	}

}

func (r *RepositoryBranch) ToProto() *v1.RepositoryBranch {
	return &v1.RepositoryBranch{
		Name:    r.Name,
		Version: r.Version,
	}
}

func RepositoryFromProto(p *v1.Repository) Repository {
	branches := make([]RepositoryBranch, len(p.GetBranches()))
	for i, branch := range p.GetBranches() {
		branches[i] = RepositoryBranchFromProto(branch)
	}

	subRepoMap := make(map[string]*Repository, len(p.GetSubRepoMap()))
	for name, repo := range p.GetSubRepoMap() {
		r := RepositoryFromProto(repo)
		subRepoMap[name] = &r
	}

	fileTombstones := make(map[string]struct{}, len(p.GetFileTombstones()))
	for _, file := range p.GetFileTombstones() {
		fileTombstones[file] = struct{}{}
	}

	return Repository{
		ID:                   p.GetId(),
		Name:                 p.GetName(),
		URL:                  p.GetUrl(),
		Source:               p.GetSource(),
		Branches:             branches,
		SubRepoMap:           subRepoMap,
		CommitURLTemplate:    p.GetCommitUrlTemplate(),
		FileURLTemplate:      p.GetFileUrlTemplate(),
		LineFragmentTemplate: p.GetLineFragmentTemplate(),
		priority:             p.GetPriority(),
		RawConfig:            p.GetRawConfig(),
		Rank:                 uint16(p.GetRank()),
		IndexOptions:         p.GetIndexOptions(),
		HasSymbols:           p.GetHasSymbols(),
		Tombstone:            p.GetTombstone(),
		LatestCommitDate:     p.GetLatestCommitDate().AsTime(),
		FileTombstones:       fileTombstones,
	}
}

func (r *Repository) ToProto() *v1.Repository {
	if r == nil {
		return nil
	}

	branches := make([]*v1.RepositoryBranch, len(r.Branches))
	for i, branch := range r.Branches {
		branches[i] = branch.ToProto()
	}

	subRepoMap := make(map[string]*v1.Repository, len(r.SubRepoMap))
	for name, repo := range r.SubRepoMap {
		subRepoMap[name] = repo.ToProto()
	}

	fileTombstones := make([]string, 0, len(r.FileTombstones))
	for file := range r.FileTombstones {
		fileTombstones = append(fileTombstones, file)
	}

	return &v1.Repository{
		Id:                   r.ID,
		Name:                 r.Name,
		Url:                  r.URL,
		Source:               r.Source,
		Branches:             branches,
		SubRepoMap:           subRepoMap,
		CommitUrlTemplate:    r.CommitURLTemplate,
		FileUrlTemplate:      r.FileURLTemplate,
		LineFragmentTemplate: r.LineFragmentTemplate,
		Priority:             r.priority,
		RawConfig:            r.RawConfig,
		Rank:                 uint32(r.Rank),
		IndexOptions:         r.IndexOptions,
		HasSymbols:           r.HasSymbols,
		Tombstone:            r.Tombstone,
		LatestCommitDate:     timestamppb.New(r.LatestCommitDate),
		FileTombstones:       fileTombstones,
	}
}

func IndexMetadataFromProto(p *v1.IndexMetadata) IndexMetadata {
	languageMap := make(map[string]uint16, len(p.GetLanguageMap()))
	for language, id := range p.GetLanguageMap() {
		languageMap[language] = uint16(id)
	}

	return IndexMetadata{
		IndexFormatVersion:    int(p.GetIndexFormatVersion()),
		IndexFeatureVersion:   int(p.GetIndexFeatureVersion()),
		IndexMinReaderVersion: int(p.GetIndexMinReaderVersion()),
		IndexTime:             p.GetIndexTime().AsTime(),
		PlainASCII:            p.GetPlainAscii(),
		LanguageMap:           languageMap,
		ZoektVersion:          p.GetZoektVersion(),
		ID:                    p.GetId(),
	}
}

func (m *IndexMetadata) ToProto() *v1.IndexMetadata {
	if m == nil {
		return nil
	}

	languageMap := make(map[string]uint32, len(m.LanguageMap))
	for language, id := range m.LanguageMap {
		languageMap[language] = uint32(id)
	}

	return &v1.IndexMetadata{
		IndexFormatVersion:    int64(m.IndexFormatVersion),
		IndexFeatureVersion:   int64(m.IndexFeatureVersion),
		IndexMinReaderVersion: int64(m.IndexMinReaderVersion),
		IndexTime:             timestamppb.New(m.IndexTime),
		PlainAscii:            m.PlainASCII,
		LanguageMap:           languageMap,
		ZoektVersion:          m.ZoektVersion,
		Id:                    m.ID,
	}
}

func RepoStatsFromProto(p *v1.RepoStats) RepoStats {
	return RepoStats{
		Repos:                      int(p.GetRepos()),
		Shards:                     int(p.GetShards()),
		Documents:                  int(p.GetDocuments()),
		IndexBytes:                 p.GetIndexBytes(),
		ContentBytes:               p.GetContentBytes(),
		NewLinesCount:              p.GetNewLinesCount(),
		DefaultBranchNewLinesCount: p.GetDefaultBranchNewLinesCount(),
		OtherBranchesNewLinesCount: p.GetOtherBranchesNewLinesCount(),
	}
}

func (s *RepoStats) ToProto() *v1.RepoStats {
	return &v1.RepoStats{
		Repos:                      int64(s.Repos),
		Shards:                     int64(s.Shards),
		Documents:                  int64(s.Documents),
		IndexBytes:                 s.IndexBytes,
		ContentBytes:               s.ContentBytes,
		NewLinesCount:              s.NewLinesCount,
		DefaultBranchNewLinesCount: s.DefaultBranchNewLinesCount,
		OtherBranchesNewLinesCount: s.OtherBranchesNewLinesCount,
	}
}

func RepoListEntryFromProto(p *v1.RepoListEntry) *RepoListEntry {
	if p == nil {
		return nil
	}

	return &RepoListEntry{
		Repository:    RepositoryFromProto(p.GetRepository()),
		IndexMetadata: IndexMetadataFromProto(p.GetIndexMetadata()),
		Stats:         RepoStatsFromProto(p.GetStats()),
	}
}

func (r *RepoListEntry) ToProto() *v1.RepoListEntry {
	if r == nil {
		return nil
	}

	return &v1.RepoListEntry{
		Repository:    r.Repository.ToProto(),
		IndexMetadata: r.IndexMetadata.ToProto(),
		Stats:         r.Stats.ToProto(),
	}
}

func MinimalRepoListEntryFromProto(p *v1.MinimalRepoListEntry) MinimalRepoListEntry {
	branches := make([]RepositoryBranch, len(p.GetBranches()))
	for i, branch := range p.GetBranches() {
		branches[i] = RepositoryBranchFromProto(branch)
	}

	return MinimalRepoListEntry{
		HasSymbols: p.GetHasSymbols(),
		Branches:   branches,
	}
}

func (m *MinimalRepoListEntry) ToProto() *v1.MinimalRepoListEntry {
	branches := make([]*v1.RepositoryBranch, len(m.Branches))
	for i, branch := range m.Branches {
		branches[i] = branch.ToProto()
	}
	return &v1.MinimalRepoListEntry{
		HasSymbols: m.HasSymbols,
		Branches:   branches,
	}
}

func RepoListFromProto(p *v1.ListResponse) *RepoList {
	repos := make([]*RepoListEntry, len(p.GetRepos()))
	for i, repo := range p.GetRepos() {
		repos[i] = RepoListEntryFromProto(repo)
	}

	reposMap := make(map[uint32]MinimalRepoListEntry, len(p.GetReposMap()))
	for id, mle := range p.GetReposMap() {
		reposMap[id] = MinimalRepoListEntryFromProto(mle)
	}

	minimal := make(map[uint32]*MinimalRepoListEntry, len(p.GetMinimal()))
	for id, mle := range p.GetMinimal() {
		m := MinimalRepoListEntryFromProto(mle)
		minimal[id] = &m
	}

	return &RepoList{
		Repos:    repos,
		ReposMap: reposMap,
		Crashes:  int(p.GetCrashes()),
		Stats:    RepoStatsFromProto(p.GetStats()),
		Minimal:  minimal,
	}
}

func (r *RepoList) ToProto() *v1.ListResponse {
	repos := make([]*v1.RepoListEntry, len(r.Repos))
	for i, repo := range r.Repos {
		repos[i] = repo.ToProto()
	}

	reposMap := make(map[uint32]*v1.MinimalRepoListEntry, len(r.ReposMap))
	for id, repo := range r.ReposMap {
		reposMap[id] = repo.ToProto()
	}

	minimal := make(map[uint32]*v1.MinimalRepoListEntry, len(r.Minimal))
	for id, repo := range r.Minimal {
		minimal[id] = repo.ToProto()
	}

	return &v1.ListResponse{
		Repos:    []*v1.RepoListEntry{},
		ReposMap: reposMap,
		Crashes:  int64(r.Crashes),
		Stats:    r.Stats.ToProto(),
		Minimal:  minimal,
	}
}

func (l *ListOptions) ToProto() *v1.ListOptions {
	if l == nil {
		return nil
	}
	var field v1.ListOptions_RepoListField
	switch l.Field {
	case RepoListFieldRepos:
		field = v1.ListOptions_REPO_LIST_FIELD_REPOS
	case RepoListFieldMinimal:
		field = v1.ListOptions_REPO_LIST_FIELD_MINIMAL
	case RepoListFieldReposMap:
		field = v1.ListOptions_REPO_LIST_FIELD_REPOS_MAP
	}

	return &v1.ListOptions{
		Field:   field,
		Minimal: l.Minimal,
	}
}

func ListOptionsFromProto(p *v1.ListOptions) *ListOptions {
	if p == nil {
		return nil
	}
	var field RepoListField
	switch p.GetField() {
	case v1.ListOptions_REPO_LIST_FIELD_REPOS:
		field = RepoListFieldRepos
	case v1.ListOptions_REPO_LIST_FIELD_MINIMAL:
		field = RepoListFieldMinimal
	case v1.ListOptions_REPO_LIST_FIELD_REPOS_MAP:
		field = RepoListFieldReposMap
	}
	return &ListOptions{
		Field:   field,
		Minimal: p.GetMinimal(),
	}
}

func SearchOptionsFromProto(p *v1.SearchOptions) *SearchOptions {
	if p == nil {
		return nil
	}

	return &SearchOptions{
		EstimateDocCount:       p.GetEstimateDocCount(),
		Whole:                  p.GetWhole(),
		ShardMaxMatchCount:     int(p.GetShardMaxMatchCount()),
		TotalMaxMatchCount:     int(p.GetTotalMaxMatchCount()),
		ShardRepoMaxMatchCount: int(p.GetShardRepoMaxMatchCount()),
		MaxWallTime:            p.GetMaxWallTime().AsDuration(),
		FlushWallTime:          p.GetFlushWallTime().AsDuration(),
		MaxDocDisplayCount:     int(p.GetMaxDocDisplayCount()),
		NumContextLines:        int(p.GetNumContextLines()),
		ChunkMatches:           p.GetChunkMatches(),
		UseDocumentRanks:       p.GetUseDocumentRanks(),
		DocumentRanksWeight:    p.GetDocumentRanksWeight(),
		Trace:                  p.GetTrace(),
		SpanContext:            p.GetSpanContext(),
	}
}

func (s *SearchOptions) ToProto() *v1.SearchOptions {
	if s == nil {
		return nil
	}

	return &v1.SearchOptions{
		EstimateDocCount:       s.EstimateDocCount,
		Whole:                  s.Whole,
		ShardMaxMatchCount:     int64(s.ShardMaxMatchCount),
		TotalMaxMatchCount:     int64(s.TotalMaxMatchCount),
		ShardRepoMaxMatchCount: int64(s.ShardRepoMaxMatchCount),
		MaxWallTime:            durationpb.New(s.MaxWallTime),
		FlushWallTime:          durationpb.New(s.FlushWallTime),
		MaxDocDisplayCount:     int64(s.MaxDocDisplayCount),
		NumContextLines:        int64(s.NumContextLines),
		ChunkMatches:           s.ChunkMatches,
		UseDocumentRanks:       s.UseDocumentRanks,
		DocumentRanksWeight:    s.DocumentRanksWeight,
		Trace:                  s.Trace,
		DebugScore:             s.DebugScore,
		SpanContext:            s.SpanContext,
	}
}
