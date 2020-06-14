// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package service

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson727fe99aDecodeGithubComNickeskovDbForumInternalPkgModelsService(in *jlexer.Lexer, out *Statuses) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Statuses, 0, 2)
			} else {
				*out = Statuses{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Status
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson727fe99aEncodeGithubComNickeskovDbForumInternalPkgModelsService(out *jwriter.Writer, in Statuses) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Statuses) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson727fe99aEncodeGithubComNickeskovDbForumInternalPkgModelsService(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Statuses) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson727fe99aEncodeGithubComNickeskovDbForumInternalPkgModelsService(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Statuses) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson727fe99aDecodeGithubComNickeskovDbForumInternalPkgModelsService(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Statuses) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson727fe99aDecodeGithubComNickeskovDbForumInternalPkgModelsService(l, v)
}
func easyjson727fe99aDecodeGithubComNickeskovDbForumInternalPkgModelsService1(in *jlexer.Lexer, out *Status) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "forum":
			out.Forum = int64(in.Int64())
		case "post":
			out.Post = int64(in.Int64())
		case "thread":
			out.Thread = int32(in.Int32())
		case "user":
			out.User = int32(in.Int32())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson727fe99aEncodeGithubComNickeskovDbForumInternalPkgModelsService1(out *jwriter.Writer, in Status) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Forum))
	}
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix)
		out.Int64(int64(in.Post))
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int32(int32(in.Thread))
	}
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		out.Int32(int32(in.User))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Status) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson727fe99aEncodeGithubComNickeskovDbForumInternalPkgModelsService1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Status) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson727fe99aEncodeGithubComNickeskovDbForumInternalPkgModelsService1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Status) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson727fe99aDecodeGithubComNickeskovDbForumInternalPkgModelsService1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Status) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson727fe99aDecodeGithubComNickeskovDbForumInternalPkgModelsService1(l, v)
}
