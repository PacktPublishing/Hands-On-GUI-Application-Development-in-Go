package main

/*
#cgo CFLAGS: -I/home/andy/Code/Go/src/github.com/andlabs/ui
#cgo LDFLAGS: /home/andy/Code/Go/src/github.com/andlabs/ui/libui_linux_amd64.a -lm -ldl
#cgo pkg-config: gtk+-3.0

#include "ui.h"

typedef struct uiWindow uiWindow;

typedef struct uiMenuItem uiMenuItem;
typedef struct uiMenu uiMenu;

void uiMenuItemOnClicked(uiMenuItem *m, void (*f)(uiMenuItem *sender, uiWindow *window, void *data), void *data);

uiMenuItem *uiMenuAppendItem(uiMenu *m, const char *name);
void uiMenuAppendSeparator(uiMenu *m);
uiMenuItem *uiMenuAppendQuitItem(uiMenu *m);

uiMenu *uiNewMenu(const char *name);

void onMenuNewClicked(uiMenuItem *item, uiWindow *w, void *data) {
	void menuNewClicked(void);
	menuNewClicked();
}

int onQuit(void *data) {
	return 1;
}

void loadMenu() {
	uiMenu *menu;
	uiMenuItem *item;

	menu = uiNewMenu("File");
	item = uiMenuAppendItem(menu, "New");
	uiMenuItemOnClicked(item, onMenuNewClicked, NULL);
	uiMenuAppendSeparator(menu);
	item = uiMenuAppendQuitItem(menu);
	uiOnShouldQuit(onQuit, NULL);

	menu = uiNewMenu("Help");
	item = uiMenuAppendItem(menu, "About");
}
*/
import "C"
