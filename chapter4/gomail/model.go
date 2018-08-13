// +build windows

package main

import "github.com/lxn/walk"

import "github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"

type EmailClientModel struct {
	Server *client.EmailServer

	root walk.TreeItem

	itemsResetPublisher  walk.TreeItemEventPublisher
	itemChangedPublisher walk.TreeItemEventPublisher
}

func (e *EmailClientModel) LazyPopulation() bool {
	return false
}

func (e *EmailClientModel) RootCount() int {
	return 1
}

func (e *EmailClientModel) RootAt(index int) walk.TreeItem {
	return e.root
}

func (e *EmailClientModel) ItemsReset() *walk.TreeItemEvent {
	return e.itemsResetPublisher.Event()
}

func (e *EmailClientModel) ItemChanged() *walk.TreeItemEvent {
	return e.itemChangedPublisher.Event()
}

func (e *EmailClientModel) SetServer(s *client.EmailServer) {
	e.Server = s

	e.root = NewInboxList(s.ListMessages())
	e.itemsResetPublisher.Publish(e.root)
}

func NewEmailClientModel() *EmailClientModel {
	return &EmailClientModel{}
}

type InboxList struct {
	emails []walk.TreeItem
}

func (i *InboxList) Text() string {
	return "Inbox"
}

func (i *InboxList) Parent() walk.TreeItem {
	return nil
}

func (i *InboxList) ChildCount() int {
	return len(i.emails)
}

func (i *InboxList) ChildAt(index int) walk.TreeItem {
	return i.emails[index]
}

func NewInboxList(l []*client.EmailMessage) *InboxList {
	list := &InboxList{}

	for _, item := range l {
		list.emails = append(list.emails, &EmailModel{item, list})
	}

	return list
}

type EmailModel struct {
	email  *client.EmailMessage
	parent walk.TreeItem
}

func (e *EmailModel) Text() string {
	return e.email.Subject
}

func (e *EmailModel) Parent() walk.TreeItem {
	return e.parent
}

func (e *EmailModel) ChildCount() int {
	return 0
}

func (e *EmailModel) ChildAt(index int) walk.TreeItem {
	return nil
}
