// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_apt.go
// -*- mode: go; coding: utf-8; -*-
// Created on 21. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-22 23:09:08 krylon>

package backend

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"regexp"
	"time"

	"github.com/blicero/krylib"
	"github.com/blicero/pkman/common"
	"github.com/blicero/pkman/database"
	"github.com/blicero/pkman/logdomain"
)

const (
	cmdApt = "/usr/bin/apt"
)

// PkgApt implements the PkgManager interface for Debian's apt.
type PkgApt struct {
	log *log.Logger
	db  *database.Database
}

// CreatePkgApt creates a PkgApt instance to interface with the apt
// package manager used by Debian, Ubuntu, and related distros.
func CreatePkgApt() (*PkgApt, error) {
	var (
		err error
		pk  = new(PkgApt)
	)

	if pk.log, err = common.GetLogger(logdomain.PkgManager); err != nil {
		return nil, err
	} else if pk.db, err = database.OpenDB(common.DbPath); err != nil {
		pk.log.Printf("[ERROR] Cannot open database at %s: %s\n",
			common.DbPath,
			err.Error())
		return nil, err
	}

	return pk, nil
} // func CreatePkgApt() (*PkgApt, error)

var patSearch = regexp.MustCompile(`^(\S+) - (.*)`) // nolint: unused

/*
Output of apt-cache search emacs
acl2-emacs - Rechenbetonte Logik für applikatives Common Lisp: Emacs-Schnittstelle
agda - Abhängig typisierte, funktionale Programmiersprache
agda-bin - Befehlszeilenschnittstelle zu Agda
agda-stdlib - Standardbibliothek für Agda
agda-stdlib-doc - Standardbibliothek für Agda - Dokumentation
alot - Textbasiertes Mailprogramm, verwendet Notmuch
alot-doc - Textbasiertes Mailprogramm, verwendet Notmuch - Dokumentation
anthy-el - Japanische Kana-Kanji-Konvertierung - Elisp-Frontend
aplus-fsf-el - XEmacs-Lisp für Entwicklungen in A+
auctex - Integrierte Umgebung für die Bearbeitung von Dokumenten mit TeX etc.
chktex - Findet typographische Fehler in LaTeX
clang-format-11 - Werkzeug zum Formatieren von C-/C++-/Obj-C-Code
clang-format-9 - Werkzeug zum Formatieren von C-/C++-/Obj-C-Code
cmucl-source - CMUCL-Lisp-Quelltext
colordiff - Werkzeug zur farbigen Hervorhebung von »diff«-Ausgaben
coq - Beweis-Assistent für Logik höherer Ordnung (Toplevel und Compiler)
crypt++el - Emacs-Lisp-Code zur Bearbeitung gepackter und verschlüsselter Dateien
cxref-emacs - Erzeugt LaTeX- und HTML-Dokumentation für C-Programme
devhelp - Hilfsprogramm für GNOME-Entwickler
dh-elpa-helper - Hilfspaket für Emacs-Lisp-Erweiterungen
dictem - Dict-Client für Emacs
dictionaries-common - Rechtschreib-Wörterbücher - gemeinsame Hilfsprogramme
eblook - Suchbefehl für elektronische Wörterbücher mit Hilfe der EB-Bibliothek
ecasound-el - Mehrspurfähiger Audiorecorder und Effektprozessor (Emacs)
eflite - Ein auf Festival-Lite basierender Sprachserver für Emacspeak
eldav - Emacs-Schnittstelle zu WebDAV-Servern
elpa-agda2-mode - Abhängig typisierte, funktionale Programmiersprache - Emacs-Modus
elpa-apache-mode - Emacs-Hauptmodus zum Bearbeiten von Apache-Konfigurationsdateien
elpa-c-sig - Signatur-Werkzeug für GNU Emacs
elpa-darcsum - PCL-CVS ähnelnde Schnittstelle zur Verwaltung von darcs-Patches
elpa-devscripts - Emacs-Ummantelung für die Befehle in devscripts
elpa-discover-my-major - Feststellen der Tastenkombinationen und ihrer Bedeutung für den aktuellen Emacs-Hauptmodus
elpa-ess - Emacs-Modus für statistische Programmierung und Datenanalyse
elpa-git-annex - Emacs-Integration für git-annex
elpa-haskell-mode - Emacs-Hauptmodus zum Editieren von Haskell
elpa-magit - Emacs-Schnittstelle für Git
elpa-protobuf-mode - Emacs-Addon zur Bearbeitung von Protocol Buffers
elpa-restart-emacs - Emacs aus Emacs heraus neu starten
elpa-verbiste - Konjugationssystem für Französisch und Italienisch - Emacs-Erweiterung
elscreen - Screen für Emacse
emacs - Editor GNU Emacs (Metapaket)
emacs-bin-common - Editor GNU Emacs - gemeinsame, architekturabhängige Dateien
emacs-common - Editor GNU Emacs - gemeinsame, architekturunabhängige Infrastruktur
emacs-el - LISP-Dateien (.el) für den Editor GNU Emacs
emacs-gtk - Editor GNU Emacs (mit Unterstützung für eine GTK+-Benutzeroberfläche)
emacs-nox - Editor GNU Emacs (ohne Unterstützung einer grafischen Oberfläche)
emacsen-common - Gemeinsame Funktionen aller Emacs-Varianten
emacspeak - Sprachausgabe-Schnittstelle für Emacs
emacspeak-ss - Emacspeak-Sprachserver für verschiedene Synthesizer
emms - Emacs-Multimediasystem
erlang - Simultane, verteilte und funktionelle Echtzeitsprache
erlang-mode - Emacs-Haupteditiermodus für Erlang
erlang-tools - Verschiedene Werkzeuge für Erlang/OTP
ess - Übergangspaket für den Wechsel von ess zu elpa-ess
etktab - ASCII Gitarrengriffschrift-Editor
exuberant-ctags - Erzeugt Indexdateien von Quelltextdefinitionen
fetchmail - SSL-fähiger E-Mail-Sammler/-Versender für POP3, APOP und IMAP
gettext-el - Emacs-Modus zur Bearbeitung der gettext-PO-Dateien
git-el - Schnelles, skalierbares, verteiltes Versionskontrollsystem (Emacs-Unterstützung)
global - Werkzeuge zum Suchen und Browsen in Quelltext
gmult - Finden Sie heraus, welche Buchstaben welche Ziffern darstellen!
gnuserv - Erlaubt das Anbinden an einen bereits laufenden Emacs
id-utils - Schnelles Werkzeug mit hohem Durchsatz für Bezeichnerdatenbanken
idl-font-lock-el - OMG IDL Schriften-Sperrung für Emacs
idn - Befehlszeilen- und Emacs-Schnittstelle für GNU Libidn
ilisp - Emacs-Schnittstelle zu LISP-Implementationen
ilisp-doc - Dokumentation für das Paket ILISP
info2man - Konvertiert GNU-Info-Dateien in POD oder Handbuchseiten
initz - Nutzung von verschieden Initdateien für Emacsen
ispell - Internationales Ispell (eine interaktive Schreibkorrektur)
jed - Editor für Programmierer (Textmodus-Version)
joe - benutzerfreundlicher Vollbild-Texteditor
jove - Jonathans Version von Emacs - ein kompakter, mächtiger Editor
kdesdk-scripts - Skripte und Datendateien für die Entwicklung
latex-cjk-common - LaTeX-Makropaket für CJK (Chinesisch/Japanisch/Koreanisch)
ledit - Zeileneditor für interaktive Programme
libghc-agda-dev - Abhängig typisierte, funktionale Programmiersprache
libghc-agda-doc - Abhängig typisierte, funktionale Programmiersprache - Dokumentation
librep-dev - Entwicklungsbibliotheken und Header für librep
liece-dcc - DCC-Programm für liece
lsdb - die Lovely Sister Database (email rolodex) für Emacs
lyskom-elisp-client - Emacs-CLient für LysKOM
malaga-mode - System für automatische Sprachanalyse - Emacs-Modus
maxima-emacs - Computeralgebrasystem -- Emacs-Schnittstelle
mew - Mailreader mit PGP/MIME-Unterstützung für Emacs
mg - Mikroskopischer Editor im Stil von GNU Emacs
mgp - MagicPoint — ein X11-basiertes Präsentationsprogramm
midge - Ein Text-zu-MIDI-Programm
mksh - MirBSD Korn Shell
mmm-mode - Mehrfacher »Major Mode« für Emacs
mpqc-support - Massiv-paralleles Quantenchemieprogramm (Hilfsprogramme)
mu4e - E-Mail-Client für Emacs, basiert auf mu (maildir-utils)
nmh - Programme für die Verarbeitung von E-Mails
nomarch - Entpackt .ARC- und .ARK-MS-DOS-Archive
post-el - Emacs-Mode zur Mailbearbeitung
proofgeneral - Generisches Frontend für Beweisassistenten
proofgeneral-doc - Generisches Frontend für Beweisassistenten - Dokumentation
pylint - Statisches Prüfprogramm für Python-3-Code und Generator von UML-Diagrammen
pymacs - Schnittstelle zwischen Emacs Lisp und Python
quilt-el - Einfache Emacs-Schnittstelle zu quilt
rdtool-elisp - Emacs-lisp rd-Modus zum Erstellen von RD-Dokumenten
remembrance-agent - Emacs-Modus zum Finden relevanter Texte
rep - Befehlsinterpreter für Lisp
rep-doc - Dokumentation für den Lisp-Befehlsinterpreter
sawfish - X11-Fenstermanager
search-ccsb - Suchprogramm für BibTeX
speechd-el-doc-cs - speechd-el Dokumentation in Tschechisch
stumpwm - Kachelnder und tastaturgesteuerter Fenstermanager geschrieben in Common Lisp
sylpheed - Schlanker E-Mail-Client mit GTK+
tdiary-mode - Emacs-Modus für Editieren von tDiary
timidity-el - Emacs-Frontend für TiMidity++
tkcon - Erweiterte interaktive Konsole für die Entwicklung mit Tcl
tmux - Terminal-Multiplexer
txt2regex - Ein Suchmuster-Zauberer auf Basis von bash2
tzc - Einfacher Zephyr-Client
uim-el - Universal Input Method - Emacs-Frontend
vile - VI wie Emacs - arbeitet wie vi
vile-common - VI Like Emacs - Support-Dateien für vile/xvile
vile-filters - VI Like Emacs - Hervorhebungsfilter für vile/xvile
xemacs21-basesupport - Editor und »Spülbecken« -- compilierte Elisp-Hilfsdateien
xemacs21-basesupport-el - Editor und Abfluss -- Elisp-Unterstützung (Quell-Dateien)
xemacs21-bin - sehr flexibler Texteditor -- benötigte Binärdateien
xemacs21-mule - sehr flexibler Texteditor -- Mule Binärdatei
xemacs21-mule-canna-wnn - Sehr flexibler Texteditor -- Mule-Binärdatei mit Unterstützung für Canna und Wnn
xemacs21-mulesupport - Editor und Küchenspüle -- Mule elisp-Unterstützungs-Dateien
xemacs21-mulesupport-el - Editor und Abfluss -- Elisp-Unterstützung (Quell-Dateien)
xemacs21-nomule - Sehr flexibler Texteditor -- nicht-Mule Binärdatei
xemacs21-support - Sehr flexibler Texteditor -- architekturunabhängige Unterstützungsdateien
xfonts-kapl - APL-Zeichensätze für A+-Entwicklung
xjed - Ein Editor für Programmierer (X11‐Version)
xvile - VI Like Emacs - VI-artiger Editor (X11)
yatex - Ein weiterer TeX-Modus für Emacs
yc-el - Noch ein weiterer Canna-Client für Emacsen
zile - Editor mit sehr kleiner Emacs-Untermenge
elpa-aggressive-indent - Emacs minor mode that reindents code after every change
apel - portable library for emacsen
elpa-assess - test support functions for Emacs
elpa-atomic-chrome - edit a web-browser text entry area with Emacs
elpa-ats2-mode - ATS version 2 programming language emacs mode
elpa-auto-complete - intelligent auto-completion extension for GNU Emacs
elpa-auto-dictionary - automatic dictionary switcher for Emacs spell checking
auto-install-el - Auto install elisp file
autodep8 - DEP-8 test control file generator
elpa-avy - jump to things in Emacs tree-style
elpa-bar-cursor - switch Emacs block cursor to a bar
bbdb - The Insidious Big Brother Database (email rolodex) for Emacs
bbdb3 - Reboot of the BBDB Insidious Big Brother Database for Emacs
elpa-bm - visual bookmarks for GNU Emacs
elpa-bongo - buffer-oriented media player for GNU Emacs
elpa-boxquote - quote text in Emacs with a semi-box
c-sig - Transition package, c-sig to elpa-c-sig
cafeobj-mode - Emacs major mode for editing CafeOBJ source code
elpa-caml - emacs mode for editing OCaml programs
elpa-ps-ccrypt - Emacs addon for working with files encrypted with ccrypt
elpa-char-menu - create your own menu for fast insertion of arbitrary symbols
cider-doc - Clojure IDE for Emacs - documentation
elpa-cider - Clojure IDE for Emacs
elpa-circe - client for IRC in Emacs
cl-iterate - Jonathan Amsterdam's Common Lisp iterator/gatherer/accumulator facility
elpa-clojure-mode - Emacs major mode for Clojure code
elpa-closql - Store EIEIO objects using EmacSQL
elpa-clues-theme - cream/brown/orange color theme for Emacs
elpa-color-theme-modern - deftheme reimplementation of classic Emacs color-themes
commit-patch - utility to commit fine grained patches to source code control repositories
elpa-company-lsp - Company completion backend for emacs lsp-mode.
elpa-company - Modular in-buffer completion framework for Emacs
elpa-csv-mode - Emacs major mode for editing comma, char, and tab separated values
elpa-cycle-quotes - Emacs command to cycle between quotation marks
darcsum - Transition package, darcsum to elpa-darcsum
elpa-dash - modern list manipulation library for Emacs
elpa-dash-functional - collection of functional combinators for Emacs Lisp
ddskk - efficient and characteristic Japanese input system for Emacs
debian-el - Transition package, debian-el to elpa-debian-el
elpa-debian-el - Emacs helpers specific to Debian users
elpa-debpaste - paste.debian.net client for Emacs
elpa-deft - Emacs mode to browse, filter, and edit plain text notes
develock-el - additional font-lock keywords for the developers on Emacs
devscripts-el - Transition package, devscripts-el to elpa-devscripts
dh-elpa - Debian helper tools for packaging emacs lisp extensions
dh-make-elpa - helper for creating Debian packages from ELPA packages
elpa-dictionary - dictionary client for Emacs
elpa-diff-hl - highlight uncommitted changes using VC
elpa-diminish - hiding or abbreviation of the mode line displays of minor-modes
elpa-dired-quick-sort - persistent quick sorting of dired buffers in various ways
elpa-dired-rsync - support for rsync from Emacs dired buffers
docbook2x - Converts DocBook/XML documents into man pages and TeXinfo
elpa-dockerfile-mode - Major mode for editing Docker's Dockerfiles
dpkg-dev-el - Transition package, dpkg-dev-el to elpa-dpkg-dev-el
elpa-dpkg-dev-el - Emacs helpers specific to Debian development
elpa-dumb-jump - jump to definition for multiple languages without configuration
e-wrapper - invoke your editor, with optional file:lineno handling
e2wm - simple window manager for emacs
e3 - very small text editor
elpa-ebib - BibTeX database manager for Emacs
ecb - code browser for Emacs supporting several languages
edict-el - Emacs interface to Edict
elpa-ediprolog - Emacs Does Interactive Prolog
elpa-editorconfig - coding style indenter for all editors - Emacsen plugin
elpa-el-mock - tiny mock and stub framework for Emacs Lisp
elpa-el-x - Emacs Lisp extensions
elpa-elfeed - Emacs Atom/RSS feed reader
elpa-bug-hunter - automatically debug and bisect your init.el or .emacs file
elpa-elisp-refs - find callers of elisp functions or macros
elpa-elisp-slime-nav - Emacs extension that provide Emacs Lisp code navigation
elpa-elm-mode - Major Emacs mode for editing Elm source code
elpa-migemo - Japanese incremental search with Romaji on Emacsen
elpa-rust-mode - Major Emacs mode for editing Rust source code
elpa-transient - Emacs key and popup interface for complex keybindings
elpa-transient-doc - Emacs key and popup interface for complex keybindings
elpa-undo-tree - Emacs minor mode for handling undo history as tree
elpa-elpher - friendly gopher and gemini client
emacs-lucid - GNU Emacs editor (with Lucid GUI support)
elpa-anzu - show number of matches in mode-line while searching
elpa-async - simple library for asynchronous processing in Emacs
elpa-bind-map - bind personal keymaps in multiple locations
elpa-buttercup - behaviour-driven testing for Emacs Lisp packages
emacs-calfw - calendar framework for Emacs
emacs-calfw-howm - calendar framework for Emacs (howm add-on)
elpa-ctable - table component for Emacs Lisp
elpa-db - database interface for Emacs Lisp
elpa-deferred - simple asynchronous functions for Emacs Lisp
elpa-epc - RPC stack for Emacs Lisp
elpa-vc-fossil - Emacs VC backend for the Fossil Version Control system
emacs-goodies-el - Miscellaneous add-ons for Emacs
elpa-helm-ag - Silver Searcher integration with Emacs Helm
elpa-highlight-indentation - highlight the indentation level in Emacs buffers
elpa-htmlize - convert buffer text and decorations to HTML
elpa-counsel - collection of Ivy-enhanced versions of common Emacs commands
elpa-ivy - generic completion mechanism for Emacs
elpa-ivy-hydra - additional key bindings for Emacs Ivy
elpa-swiper - alternative to Emacs' isearch, with an overview
elpa-jabber - Jabber client for Emacsen
emacs-jabber - Transition package, emacs-jabber to elpa-jabber
elpa-jedi - Python auto-completion for Emacs
elpa-jedi-core - common code of jedi.el and company-jedi.el
elpa-kv - key/value data structure functions for Emacs Lisp
elpa-lsp-haskell - Haskell support for lsp-mode
elpa-neotree - directory tree sidebar for Emacs that is like NERDTree for Vim
elpa-noflet - Emacs Lisp noflet macro for dynamic, local advice
elpa-openwith - seamlessly open files in external programs with Emacs
elpa-orgalist - Manage Org-like lists in non-Org Emacs buffers
elpa-pdf-tools - Display and interact with pdf in Emacs
elpa-pdf-tools-server - server for Emacs's pdf-tools
elpa-pg - Emacs Lisp interface for PostgreSQL
elpa-pod-mode - Emacs major mode for editing .pod files
elpa-powerline - Emacs version of the Vim powerline
elpa-python-environment - virtualenv API for Emacs Lisp
elpa-session - use variables, registers and buffer places across sessions
elpa-smeargle - highlight region by last updated time
elpa-uuid - UUID/GUID library for Emacs Lisp
elpa-websocket - Emacs WebSocket client and server
elpa-wgrep - edit multiple Emacs buffers using a master grep pattern buffer
elpa-wgrep-ack - edit multiple Emacs buffers using a master ack pattern buffer
elpa-wgrep-ag - edit multiple Emacs buffers using a master ag pattern buffer
elpa-wgrep-helm - edit multiple Emacs buffers with a helm-grep-mode buffer
elpa-which-key - display available keybindings in popup
emacs-window-layout - window layout manager for emacs
elpa-world-time-mode - Emacs mode to compare timezones throughout the day
emacspeak-espeak-server - espeak synthesis server for emacspeak
elpa-emacsql - high level SQL database frontend for Emacs
elpa-emacsql-mysql - high level SQL database frontend for Emacs
elpa-emacsql-psql - high level SQL database frontend for Emacs
elpa-emacsql-sqlite - high level SQL database frontend for Emacs
elpa-emacsql-sqlite3 - Yet another EmacSQL backend for SQLite
elpa-engine-mode - define and query search engines from within Emacs
elpa-epl - Emacs Package Library
elpa-eproject - assign files to Emacs projects, programmatically
elpa-ert-async - asynchronous tests for the Emacs ERT testing framework
elpa-ert-expectations - very simple unit test framework for Emacs Lisp
elpa-eshell-git-prompt - Eshell prompt themes for Git users
elpa-eshell-z - cd to frequent directory in eshell
elpa-esup - Emacs StartUp Profiler
elpa-esxml - XML, ESXML and SXML library for Emacs Lisp
elpa-evil - extensible vi layer for Emacs
elpa-evil-paredit - emacs extension, integrating evil and paredit
eweouz - Emacs interface to Evolution Data Server
elpa-exec-path-from-shell - get environment variables such as $PATH from the shell
elpa-expand-region - Increase selected region in Emacs by semantic units
elpa-eyebrowse - simple-minded way of managing window configs in Emacs
elpa-f - modern API for working with files and directories in Emacs Lisp
elpa-faceup - Regression test system for font-lock
elpa-fill-column-indicator - graphically indicate the fill column
elpa-find-file-in-project - quick access to project files in Emacs
findent - indents/converts Fortran sources
flim - library about internet message for emacsen
elpa-flx - sorting algorithm for fuzzy matching in Emacs
elpa-flx-ido - allows Emacs Ido to use the flx sorting algorithm
elpa-flycheck - modern on-the-fly syntax checking for Emacs
flycheck-doc - modern on-the-fly syntax checking for Emacs - documentation
elpa-flycheck-package - flycheck checker for Elisp package authors
elpa-folding - folding-editor minor mode for Emacs
fortran-language-server - Fortran Language Server for the Language Server Protocol
elpa-fountain-mode - Emacs major mode for screenwriting in Fountain markup
elpa-fricas - General purpose computer algebra system: emacs support
elpa-fsm - state machine library
elpa-geiser - enhanced Scheme interaction mode for Emacs
geiser - Transition Package, geiser to elpa-geiser
elpa-ggtags - improved Emacs interface to GNU GLOBAL
elpa-git-auto-commit-mode - Emacs Minor mode to automatically commit and push with git
elpa-git-timemachine - walk through git revisions of a file
elpa-gitlab-ci-mode - Emacs mode for editing GitLab CI files
gnu-smalltalk-el - GNU Smalltalk Emacs front-end
elpa-gnuplot-mode - Gnuplot mode for Emacs
gnuplot-mode - Transition Package, gnuplot-mode to elpa-gnuplot-mode
elpa-go-mode - Emacs mode for editing Go code
golang-mode - Emacs mode for editing Go code -- transitional package
goby - WYSIWYG presentation tool for Emacs
elpa-golden-ratio - automatic resizing of Emacs windows to the golden ratio
elpa-goo - generic object-orientator (Emacs support)
elpa-goto-chg - navigate the point to the most recent edit in the buffer
gramadoir - Irish language grammar checker (integration scripts)
elpa-graphviz-dot-mode - Emacs mode for the dot-language used by graphviz.
gri-el - Emacs major-mode for gri, a language for scientific graphics
haml-elisp - Emacs Lisp mode for the Haml markup language
libghc-yi-keymap-emacs-dev - Emacs keymap for Yi editor
libghc-yi-keymap-emacs-doc - Emacs keymap for Yi editor; documentation
libghc-yi-keymap-emacs-prof - Emacs keymap for Yi editor; profiling libraries
elpa-helm - Emacs incremental completion and selection narrowing framework
elpa-helm-core - Emacs Helm library files
elpa-helm-org - Emacs Helm for Org-mode headlines and keywords completion
howm - Note-taking tool on Emacs
elpa-ht - hash table library for Emacs
elpa-hydra - make Emacs bindings that stick around
elpa-ibuffer-projectile - group buffers in ibuffer list by Projectile project
elpa-ibuffer-vc - group ibuffer list by VC project and show VC status
elpa-ido-completing-read+ - completing-read-function using ido
elpa-iedit - edit multiple regions in the same way simultaneously
elpa-imenu-list - show the current Emacs buffer's imenu entries in a separate window
elpa-initsplit - code to split customizations into different files
inotify-hookable - blocking command-line interface to inotify
emacs-intl-fonts - fonts to allow multilingual PostScript printing from Emacs
elpa-irony - Emacs C/C++ minor mode powered by libclang
irony-server - Emacs C/C++ minor mode powered by libclang (server)
jedit - Plugin-based editor for programmers
elpa-jinja2-mode - Emacs major mode for editing jinja2 code
libjline2-java - console input handling in Java
elpa-js2-mode - Emacs mode for editing Javascript programs
js2-mode - Emacs mode for editing Javascript programs (dummy package)
jupp - user friendly full screen text editor
elpa-key-chord - map pairs of simultaneously pressed keys to commands
elpa-lbdb - Little Brother's DataBase Emacs extensions
elpa-ledger - command-line double-entry accounting program (emacs interface)
libledit-ocaml-dev - OCaml line editor library
elpa-let-alist - let-bind values of an assoc-list by their names in Emacs Lisp
libconfig-find-perl - module to search configuration files using OS dependent heuristics
gir1.2-kkc-1.0 - GObject introspection data for libkkc
libkkc-common - Japanese Kana Kanji input library - common data
libkkc-dev - Japanese Kana Kanji input library - development files
libkkc-utils - Japanese Kana Kanji input library - testing utility
libkkc2 - Japanese Kana Kanji input library
libkkc-data - language model data for libkkc
liblatex-table-perl - Perl extension for the automatic generation of LaTeX tables
libparse-exuberantctags-perl - exuberant ctags parser for Perl
libproc-invokeeditor-perl - Perl extension for starting a text editor
librep16 - embedded lisp command interpreter library
libtext-findindent-perl - module to heuristically determine indentation style
liece - IRC (Internet Relay Chat) client for Emacs
elpa-linum-relative - display relative line number in Emacs
liquidsoap-mode - Emacs mode for editing Liquidsoap code
lisaac-mode - Emacs mode for editing Lisaac programs
clang-format-13 - Tool to format C/C++/Obj-C code
elpa-load-relative - relative file load (within a multi-file Emacs package)
lookup-el - emacsen interface to electronic dictionaries
elpa-loop - friendly imperative loop structures for Emacs Lisp
elpa-lsp-mode - Emacs client/library for the Language Server Protocol
elpa-lua-mode - Emacs major-mode for editing Lua programs
elpa-macaulay2 - Software system for algebraic geometry research (Emacs package)
elpa-git-commit - Major mode for editing git commit message
elpa-magit-forge - Work with Git forges from the comfort of Magit
elpa-magit-popup - Use popup like Magit
elpa-mailscripts - Emacs functions for accessing tools in the mailscripts package
elpa-makey - flexible context menu system
elpa-markdown-mode - mode for editing Markdown-formatted text files in GNU Emacs
elpa-markdown-toc - Emacs TOC (table of contents) generator for markdown files
elpa-meson-mode - Major mode for the Meson build system files
elpa-message-templ - templates for Emacs message-mode
ocaml-core - OCaml core tools (metapackage)
mew-beta - mail reader supporting PGP/MIME for Emacs (development version)
mhc - schedule management tool for Emacs
mhc-utils - utilities for the MHC schedule management system
minlog - Proof assistant based on first order natural deduction calculus
mit-scheme - MIT/GNU Scheme development environment
mit-scheme-dbg - MIT/GNU Scheme debugging files
mit-scheme-doc - MIT/GNU Scheme documentation
elpa-mocker - mocking framework for Emacs
elpa-modus-themes - set of accessible themes conforming with WCAG AAA accessibility standard
elpa-monokai-theme - fruity color theme for Emacs
emacs-mozc - Mozc for Emacs
emacs-mozc-bin - Helper module for emacs-mozc
mu-cite - message citation utility for emacsen
elpa-muse - author and publish projects using Wiki-like markup
elpa-mutt-alias - Emacs package to lookup and insert expanded Mutt mail aliases
elpa-muttrc-mode - Emacs major mode for editing muttrc
nescc - Programming Language for Deeply Networked Systems
ng-cjk - Nihongo MicroGnuEmacs with CJK support
ng-cjk-canna - Nihongo MicroGnuEmacs with CJK and Canna support
ng-common - Common files used by ng-* packages
ng-latin - Nihongo MicroGnuEmacs with Latin support
elpa-no-littering - help keeping ~/.emacs.d clean
elpa-nose - easy Python test running in Emacs
elpa-notmuch - thread-based email index, search and tagging (emacs interface)
notmuch-addrlookup - Address lookup tool for Notmuch
elpa-nov - featureful EPUB (ebook) reader mode for Emacs
libre-ocaml-dev - regular expression library for OCaml
libocp-indent-ocaml - OCaml indentation tool for emacs and vim - libraries
libocp-indent-ocaml-dev - OCaml indentation tool for emacs and vim - development libraries
ocp-indent - OCaml indentation tool for emacs and vim - runtime
elpa-olivetti - Emacs minor mode to more comfortably read and write long-lined prose
oneliner-el - extensions of Emacs standard shell-mode
elpa-org-drill - emacs org-mode contrib for self-testing using spaced repetition
elpa-org - Keep notes, maintain ToDo lists, and do project planning in emacs
org-mode - Transition Package, org-mode to elpa-org
elpa-org-roam - non-hierarchical note-taking with Emacs Org-mode
org-roam-doc - non-hierarchical note-taking with Emacs Org-mode -- documentation
otags - tags file generator for OCaml
elpa-package-lint - linting library for Elisp package authors
elpa-package-lint-flymake - package-lint Flymake backend
libghc-pandoc-dev - general markup converter - libraries
libghc-pandoc-doc - general markup converter - library documentation
libghc-pandoc-prof - general markup converter - profiling libraries
pandoc - general markup converter
pandoc-data - general markup converter - data files
elpa-paredit - Emacs minor mode for structurally editing Lisp code
elpa-paredit-everywhere - cut-down version of paredit for non-lisp buffers
elpa-parent-mode - get major mode's parent modes
elpa-parsebib - Emacs Lisp library for parsing .bib files
libpcre-ocaml - OCaml bindings for PCRE (runtime)
libpcre-ocaml-dev - OCaml bindings for PCRE (Perl Compatible Regular Expression)
elpa-pcre2el - Emacs mode to convert between PCRE, Emacs and rx regexp syntax
elpa-persist - persist variables between Emacs Sessions
elpa-perspective - tagged workspaces in Emacs
elpa-php-mode - PHP Mode for GNU Emacs
elpa-pip-requirements - major mode for editing pip requirements files
elpa-pkg-info - provide information about Emacs packages
elpa-popup - visual popup user interface library for Emacs
elpa-pos-tip - Show tooltip at point
elpa-project - Emacs library for operations on the current project
elpa-projectile - project interaction library for Emacs
projectile-doc - project interaction library for Emacs - documentation
psgml - Emacs major mode for editing SGML documents
elpa-puppet-mode - Emacs major mode for Puppet manifests
python3-editor - programmatically open an editor, capture the result - Python 3.x
python3-epc - RPC stack for Emacs Lisp (Python3 version)
elpa-pyvenv - Python virtual environment interface
elpa-qml-mode - Emacs major mode for editing QT Declarative (QML) code
elpa-queue - queue data structure for Emacs Lisp
r-cran-progress - GNU R terminal progress bars
rabbit-mode - Emacs-lisp rabbit-mode for writing RD document using Rabbit
elpa-racket-mode - emacs support for editing and running racket code
rail - Replace Agent-string Internal Library
elpa-rainbow-delimiters - Emacs mode to colour-code delimiters according to their depth
python3-readlike - GNU Readline-like line editing module
elpa-redtick - tiny pomodoro timer for Emacs
elpa-relint - Emacs Lisp regexp mistake finder
reposurgeon - Tool for editing version-control repository history
elpa-rich-minority - clean-up and beautify the list of minor-modes in Emacs' mode-line
riece - IRC client for Emacs
librobert-hooke-clojure - Function wrapper library for Clojure
elpa-rtags - emacs front-end for RTags
rtags - C/C++ client/server indexer with integration for Emacs
ruby-github-markup - GitHub Markup rendering
ruby-notiffany - Wrapper libray for most popular notification libraries
ruby-org - Emacs org-mode parser for Ruby
elpa-s - string manipulation library for Emacs
sass-elisp - Emacs Lisp mode for the Sass markup language
elpa-scala-mode - Emacs major mode for editing scala source code
cmuscheme48-el - Emacs mode specialized for Scheme48
search-citeseer - BibTeX search tool
select-xface - utility for selecting X-Face on emacsen
semi - library to provide MIME feature for emacsen
sepia - Simple Emacs-Perl InterAction
elpa-seq - sequence manipulation functions for Emacs Lisp
elpa-sesman - session manager for Emacs IDEs
elpa-shut-up - Emacs Lisp macros to quieten Emacs
elpa-ag - Emacs frontend to ag
singular-ui-emacs - Computer Algebra System for Polynomial Computations -- emacs user interface
sisu - documents - structuring, publishing in multiple formats and search
cl-swank - Superior Lisp Interaction Mode for Emacs (Lisp-side server)
slime - Superior Lisp Interaction Mode for Emacs (client)
elpa-smart-mode-line - powerful and beautiful mode-line for Emacs
elpa-smart-mode-line-powerline-theme - Smart Mode Line themes that use Emacs Powerline
elpa-smex - M-x interface for Emacs with Ido-style fuzzy matching
elpa-sml-mode - Emacs major mode for editing Standard ML programs
elpa-solarized-theme - port of Solarized theme to Emacs
speechd-el - Emacs speech client using Speech Dispatcher
speechd-up - Interface between Speech Dispatcher and SpeakUp
elpa-spinner - spinner for the Emacs modeline for operations in progress
stow - Organizer for /usr/local software packages
elpa-suggest - discover Emacs Lisp functions based on examples
elpa-super-save - auto-save buffers, based on your activity
supercollider-emacs - SuperCollider mode for Emacs
elpa-sxiv - run the sxiv image viewer
elpa-systemd - major mode for editing systemd units
t-code - Japanese direct input method environment for emacsen
t-code-common - Japanese direct input method environment - common files
elpa-tabbar - Emacs minor mode that displays a tab bar at the top
xfonts-thai-etl - Thai etl fonts for X
xfonts-thai-poonlap - Poonlap Veerathanabutr's bitmap fonts for X
tiarra-conf-el - edit mode for tiarra.conf
tpp - text presentation program
elpa-transmission - Emacs interface to a Transmission session
elpa-tuareg - emacs-mode for OCaml programs
tuareg-mode - transitional package, tuareg-mode to elpa-tuareg
tweak - Efficient text-mode hex editor
twittering-mode - Twitter client for Emacs
uim-latin - Universal Input Method - Latin script input support metapackage
elpa-undercover - test coverage library for Emacs Lisp
universal-ctags - build tag file indexes of source code definitions
elpa-bind-key - simple way to manage personal keybindings
elpa-use-package - configuration macro for simplifying your .emacs
libutop-ocaml - improved OCaml toplevel (runtime library)
libutop-ocaml-dev - improved OCaml toplevel (development tools)
utop - improved OCaml toplevel
elpa-vala-mode - Emacs editor major mode for vala source code
vim-voom - Vim two-pane outliner
elpa-vimish-fold - fold text in GNU Emacs like in Vim
elpa-virtualenvwrapper - featureful virtualenv tool for Emacs
elpa-visual-fill-column - Emacs mode that wraps visual-line-mode buffers at fill-column
elpa-visual-regexp - in-buffer visual feedback while using Emacs regexps
vm - mail user agent for Emacs
elpa-volume - tweak your sound card volume from Emacs
w3m-el - simple Emacs interface of w3m
w3m-el-snapshot - simple Emacs interface of w3m (development version)
elpa-wc-mode - display a word count in the Emacs modeline
elpa-web-mode - major emacs mode for editing web templates
elpa-weechat - Chat via WeeChat's relay protocol in Emacs.
whizzytex - WYSIWYG emacs environment for LaTeX
windows-el - window manager for GNU Emacs
elpa-with-editor - call program using Emacs as $EDITOR
elpa-with-simulated-input - macro to simulate user input non-interactively
wl - mail/news reader supporting IMAP for emacsen
wl-beta - mail/news reader supporting IMAP for emacsen (development version)
wordwarvi - retro-styled side-scrolling shoot'em up arcade game
wordwarvi-sound - retro-styled side-scrolling shoot'em up arcade game [Sound Files]
elpa-writegood-mode - Emacs minor mode that provides hints for common English writing problems
elpa-ws-butler - unobtrusively remove trailing whitespace in Emacs
x-face-el - utility for displaying X-Face on emacsen
elpa-xcite - exciting cite utility for Emacsen
xcite - Transition Package, xcite to elpa-xcite
elpa-xcscope - Interactively examine a C program source in emacs
xcscope-el - Transition Package, xcscope-el to elpa-xcscope
xemacs21 - highly customizable text editor metapackage
xemacs21-supportel - highly customizable text editor -- non-required library files
xfonts-terminus-oblique - Oblique version of the Terminus font
elpa-xml-rpc - Emacs Lisp XML-RPC client
elpa-xr - convert string regexp to rx notation
elpa-xref - Library for cross-referencing commands in Emacs
xstow - Extended replacement of GNU Stow
elpa-yaml-mode - Emacs major mode for YAML files
elpa-yasnippet - template system for Emacs
yasnippet - transition Package, yasnippet to elpa-yasnippet
elpa-yasnippet-snippets - Andrea Crotti's official YASnippet snippets
yasr - General-purpose console screen reader
yorick - interpreted language and scientific graphics
elpa-zenburn-theme - low contrast color theme for Emacs
elpa-ztree - text mode directory tree
elpa-elfeed-web - Emacs Atom/RSS feed reader - web interface
wnn7egg - Wnn-nana-tamago -- EGG Input Method with Wnn7 for Emacsen
emacs-common-non-dfsg - GNU Emacs common non-DFSG items, including the core documentation
org-mode-doc - keep notes, maintain ToDo lists, and do project planning in emacs

Output of apt search emacs:

WARNING: apt does not have a stable CLI interface. Use with caution in scripts.

Sortierung…
Volltextsuche…
acl2-emacs/stable,stable 8.3dfsg-2 all
  Rechenbetonte Logik für applikatives Common Lisp: Emacs-Schnittstelle

agda/stable,stable 2.6.1-1 all
  Abhängig typisierte, funktionale Programmiersprache

agda-bin/stable 2.6.1-1+b2 amd64
  Befehlszeilenschnittstelle zu Agda

agda-stdlib/stable,stable 1.3-2 all
  Standardbibliothek für Agda

agda-stdlib-doc/stable,stable 1.3-2 all
  Standardbibliothek für Agda - Dokumentation

alot/stable,stable 0.9.1-2 all
  Textbasiertes Mailprogramm, verwendet Notmuch

alot-doc/stable,stable 0.9.1-2 all
  Textbasiertes Mailprogramm, verwendet Notmuch - Dokumentation

anthy-el/stable,stable 1:0.4-2 all
  Japanische Kana-Kanji-Konvertierung - Elisp-Frontend

apel/stable,stable 10.8+0.20201106-1 all
  portable library for emacsen

aplus-fsf-el/stable,stable 4.22.1-10.1 all
  XEmacs-Lisp für Entwicklungen in A+

auctex/stable,stable 12.2-1 all
  Integrierte Umgebung für die Bearbeitung von Dokumenten mit TeX etc.

auto-install-el/stable,stable 1.58-1.1 all
  Auto install elisp file

autodep8/stable,stable 0.24 all
  DEP-8 test control file generator

bbdb/stable,stable 3.0.1 all
  The Insidious Big Brother Database (email rolodex) for Emacs

bbdb3/stable,stable 3.2-10 all
  Reboot of the BBDB Insidious Big Brother Database for Emacs

c-sig/stable,stable 3.8-24 all
  Transition package, c-sig to elpa-c-sig

cafeobj-mode/stable,stable 1.6.0-2 all
  Emacs major mode for editing CafeOBJ source code

chktex/stable 1.7.6-4 amd64
  Findet typographische Fehler in LaTeX

cider-doc/stable,stable 0.19.0+dfsg-2.1 all
  Clojure IDE for Emacs - documentation

cl-iterate/stable,stable 20180228-1.1 all
  Jonathan Amsterdam's Common Lisp iterator/gatherer/accumulator facility

cl-swank/stable,stable 2:2.26.1+dfsg-2 all
  Superior Lisp Interaction Mode for Emacs (Lisp-side server)

clang-format-11/stable 1:11.0.1-2 amd64
  Werkzeug zum Formatieren von C-/C++-/Obj-C-Code

clang-format-13/stable 1:13.0.1-6~deb11u1 amd64
  Tool to format C/C++/Obj-C code

clang-format-9/stable 1:9.0.1-16.1 amd64
  Werkzeug zum Formatieren von C-/C++-/Obj-C-Code

cmucl-source/stable,stable 21d-1.1 all
  CMUCL-Lisp-Quelltext

cmuscheme48-el/stable,stable 1.9.2-2 all
  Emacs mode specialized for Scheme48

colordiff/stable,stable 1.0.18-1.1 all
  Werkzeug zur farbigen Hervorhebung von »diff«-Ausgaben

commit-patch/stable,stable 2.6-2.1 all
  utility to commit fine grained patches to source code control repositories

coq/stable 8.12.0-3+b3 amd64
  Beweis-Assistent für Logik höherer Ordnung (Toplevel und Compiler)

crypt++el/stable,stable 2.94-3.1 all
  Emacs-Lisp-Code zur Bearbeitung gepackter und verschlüsselter Dateien

cxref-emacs/stable,stable 1.6e-3.1 all
  Erzeugt LaTeX- und HTML-Dokumentation für C-Programme

darcsum/stable,stable 1.10+20120116-4 all
  Transition package, darcsum to elpa-darcsum

ddskk/stable,stable 17.1-4+deb11u1 all
  efficient and characteristic Japanese input system for Emacs

debian-el/stable,stable 37.10 all
  Transition package, debian-el to elpa-debian-el

develock-el/stable,stable 0.47-3.1 all
  additional font-lock keywords for the developers on Emacs

devhelp/stable 3.38.1-1 amd64
  Hilfsprogramm für GNOME-Entwickler

devscripts-el/stable,stable 40.5 all
  Transition package, devscripts-el to elpa-devscripts

dh-elpa/stable,stable 2.0.8 all
  Debian helper tools for packaging emacs lisp extensions

dh-elpa-helper/stable,stable,now 2.0.8 all  [Installiert,automatisch]
  Hilfspaket für Emacs-Lisp-Erweiterungen

dh-make-elpa/stable,stable 0.19.1 all
  helper for creating Debian packages from ELPA packages

dictem/stable,stable 1.0.4-4.1 all
  Dict-Client für Emacs

dictionaries-common/stable,stable,now 1.28.4 all  [Installiert,automatisch]
  Rechtschreib-Wörterbücher - gemeinsame Hilfsprogramme

docbook2x/stable 0.8.8-17+b1 amd64
  Converts DocBook/XML documents into man pages and TeXinfo

dpkg-dev-el/stable,stable 37.9 all
  Transition package, dpkg-dev-el to elpa-dpkg-dev-el

e-wrapper/stable,stable 0.2-1 all
  invoke your editor, with optional file:lineno handling

e2wm/stable,stable 1.4-3 all
  simple window manager for emacs

e3/stable 1:2.82+dfsg-2 amd64
  very small text editor

eblook/stable 1:1.6.1-16 amd64
  Suchbefehl für elektronische Wörterbücher mit Hilfe der EB-Bibliothek

ecasound-el/stable,stable 2.9.3-2 all
  Mehrspurfähiger Audiorecorder und Effektprozessor (Emacs)

ecb/stable,stable 2.50+git20170628-1 all
  code browser for Emacs supporting several languages

edict-el/stable,stable 1.06-11.1 all
  Emacs interface to Edict

eflite/stable 0.4.1-12 amd64
  Ein auf Festival-Lite basierender Sprachserver für Emacspeak

eldav/stable,stable 0.8.1-10.1 all
  Emacs-Schnittstelle zu WebDAV-Servern

elpa-ag/stable,stable,now 0.48-1 all  [Installiert,automatisch]
  Emacs frontend to ag

elpa-agda2-mode/stable,stable 2.6.1-1 all
  Abhängig typisierte, funktionale Programmiersprache - Emacs-Modus

elpa-aggressive-indent/stable,stable 1.9.0-3 all
  Emacs minor mode that reindents code after every change

elpa-anzu/stable,stable 0.64-1 all
  show number of matches in mode-line while searching

elpa-apache-mode/stable,stable 2.2.0-3 all
  Emacs-Hauptmodus zum Bearbeiten von Apache-Konfigurationsdateien

elpa-assess/stable,stable 0.6-1 all
  test support functions for Emacs

elpa-async/stable,stable 1.9.4-2 all
  simple library for asynchronous processing in Emacs

elpa-atomic-chrome/stable,stable 2.0.0-2 all
  edit a web-browser text entry area with Emacs

elpa-ats2-mode/stable,stable 0.4.0-1 all
  ATS version 2 programming language emacs mode

elpa-auto-complete/stable,stable 1.5.1-0.2 all
  intelligent auto-completion extension for GNU Emacs

elpa-auto-dictionary/stable,stable 1.1+14.gb364e08-1 all
  automatic dictionary switcher for Emacs spell checking

elpa-avy/stable,stable 0.5.0-2 all
  jump to things in Emacs tree-style

elpa-bar-cursor/stable,stable 2.0-1.1 all
  switch Emacs block cursor to a bar

elpa-bind-key/stable,stable 2.4.1-1 all
  simple way to manage personal keybindings

elpa-bind-map/stable,stable 1.1.1-5 all
  bind personal keymaps in multiple locations

elpa-bm/stable,stable 201905-2 all
  visual bookmarks for GNU Emacs

elpa-bongo/stable,stable 1.1-2 all
  buffer-oriented media player for GNU Emacs

elpa-boxquote/stable,stable 2.2-1 all
  quote text in Emacs with a semi-box

elpa-bug-hunter/stable,stable 1.3.1+repack-5 all
  automatically debug and bisect your init.el or .emacs file

elpa-buttercup/stable,stable 1.24-1 all
  behaviour-driven testing for Emacs Lisp packages

elpa-c-sig/stable,stable 3.8-24 all
  Signatur-Werkzeug für GNU Emacs

elpa-caml/stable,stable 4.06-2 all
  emacs mode for editing OCaml programs

elpa-char-menu/stable,stable 0.1.1-3 all
  create your own menu for fast insertion of arbitrary symbols

elpa-cider/stable,stable 0.19.0+dfsg-2.1 all
  Clojure IDE for Emacs

elpa-circe/stable,stable 2.11-2 all
  client for IRC in Emacs

elpa-clojure-mode/stable,stable 5.10.0-3 all
  Emacs major mode for Clojure code

elpa-closql/stable,stable 1.0.4-2 all
  Store EIEIO objects using EmacSQL

elpa-clues-theme/stable,stable 1.0.1-2.1 all
  cream/brown/orange color theme for Emacs

elpa-color-theme-modern/stable,stable 0.0.3-1 all
  deftheme reimplementation of classic Emacs color-themes

elpa-company/stable,stable 0.9.13-2 all
  Modular in-buffer completion framework for Emacs

elpa-company-lsp/stable,stable 2.1.0-3 all
  Company completion backend for emacs lsp-mode.

elpa-counsel/stable,stable 0.13.0-1 all
  collection of Ivy-enhanced versions of common Emacs commands

elpa-csv-mode/stable,stable 1.12-1 all
  Emacs major mode for editing comma, char, and tab separated values

elpa-ctable/stable,stable 0.1.2-6 all
  table component for Emacs Lisp

elpa-cycle-quotes/stable,stable 0.1-4 all
  Emacs command to cycle between quotation marks

elpa-darcsum/stable,stable 1.10+20120116-4 all
  PCL-CVS ähnelnde Schnittstelle zur Verwaltung von darcs-Patches

elpa-dash/stable,stable,now 2.17.0+dfsg-1 all  [Installiert,automatisch]
  modern list manipulation library for Emacs

elpa-dash-functional/stable,stable,now 1.2.0+dfsg-7 all  [Installiert,automatisch]
  collection of functional combinators for Emacs Lisp

elpa-db/stable,stable 0.0.6+git20140421.b3a423f-3 all
  database interface for Emacs Lisp

elpa-debian-el/stable,stable 37.10 all
  Emacs helpers specific to Debian users

elpa-debpaste/stable,stable 0.1.5-4 all
  paste.debian.net client for Emacs

elpa-deferred/stable,stable 0.5.1-4 all
  simple asynchronous functions for Emacs Lisp

elpa-deft/stable,stable 0.8-3 all
  Emacs mode to browse, filter, and edit plain text notes

elpa-devscripts/stable,stable 40.5 all
  Emacs-Ummantelung für die Befehle in devscripts

elpa-dictionary/stable,stable 1.10+git20190107-3 all
  dictionary client for Emacs

elpa-diff-hl/stable,stable 1.8.8-1 all
  highlight uncommitted changes using VC

elpa-diminish/stable,stable 0.45-4 all
  hiding or abbreviation of the mode line displays of minor-modes

elpa-dired-quick-sort/stable,stable 0.1.1-1 all
  persistent quick sorting of dired buffers in various ways

elpa-dired-rsync/stable,stable 0.6-1 all
  support for rsync from Emacs dired buffers

elpa-discover-my-major/stable,stable 1.0-4 all
  Feststellen der Tastenkombinationen und ihrer Bedeutung für den aktuellen Emacs-Hauptmodus

elpa-dockerfile-mode/stable,stable 1.2-2 all
  Major mode for editing Docker's Dockerfiles

elpa-dpkg-dev-el/stable,stable 37.9 all
  Emacs helpers specific to Debian development

elpa-dumb-jump/stable,stable 0.5.3-1 all
  jump to definition for multiple languages without configuration

elpa-ebib/stable,stable 2.15.4-3 all
  BibTeX database manager for Emacs

elpa-ediprolog/stable,stable 2.1-1 all
  Emacs Does Interactive Prolog

elpa-editorconfig/stable,stable 0.8.1-3 all
  coding style indenter for all editors - Emacsen plugin

elpa-el-mock/stable,stable 1.25.1-4 all
  tiny mock and stub framework for Emacs Lisp

elpa-el-x/stable,stable 0.3.1-4 all
  Emacs Lisp extensions

elpa-elfeed/stable,stable 3.4.1-1 all
  Emacs Atom/RSS feed reader

elpa-elfeed-web/stable,stable 3.4.1-1 all
  Emacs Atom/RSS feed reader - web interface

elpa-elisp-refs/stable,stable 1.3-3 all
  find callers of elisp functions or macros

elpa-elisp-slime-nav/stable,stable 0.9-5 all
  Emacs extension that provide Emacs Lisp code navigation

elpa-elm-mode/stable,stable 0.20.3-3 all
  Major Emacs mode for editing Elm source code

elpa-elpher/stable,stable 2.10.2-2 all
  friendly gopher and gemini client

elpa-emacsql/stable,stable 3.0.0+ds-2 all
  high level SQL database frontend for Emacs

elpa-emacsql-mysql/stable,stable 3.0.0+ds-2 all
  high level SQL database frontend for Emacs

elpa-emacsql-psql/stable,stable 3.0.0+ds-2 all
  high level SQL database frontend for Emacs

elpa-emacsql-sqlite/stable 3.0.0+ds-2 amd64
  high level SQL database frontend for Emacs

elpa-emacsql-sqlite3/stable,stable 1.0.2-1 all
  Yet another EmacSQL backend for SQLite

elpa-engine-mode/stable,stable 2.1.1-1 all
  define and query search engines from within Emacs

elpa-epc/stable,stable 0.1.1-6 all
  RPC stack for Emacs Lisp

elpa-epl/stable,stable 0.9-3 all
  Emacs Package Library

elpa-eproject/stable,stable 1.5+git20180312.068218d-3 all
  assign files to Emacs projects, programmatically

elpa-ert-async/stable,stable 0.1.2-5 all
  asynchronous tests for the Emacs ERT testing framework

elpa-ert-expectations/stable,stable 0.2-4 all
  very simple unit test framework for Emacs Lisp

elpa-eshell-git-prompt/stable,stable 0.1.2-4 all
  Eshell prompt themes for Git users

elpa-eshell-z/stable,stable 0.4-3 all
  cd to frequent directory in eshell

elpa-ess/stable,stable 18.10.2-2 all
  Emacs-Modus für statistische Programmierung und Datenanalyse

elpa-esup/stable,stable 0.7.1-3 all
  Emacs StartUp Profiler

elpa-esxml/stable,stable 0.3.5-1 all
  XML, ESXML and SXML library for Emacs Lisp

elpa-evil/stable,stable 1.14.0-1 all
  extensible vi layer for Emacs

elpa-evil-paredit/stable,stable 0.0.2-5 all
  emacs extension, integrating evil and paredit

elpa-exec-path-from-shell/stable,stable 1.12-2 all
  get environment variables such as $PATH from the shell

elpa-expand-region/stable,stable 0.11.0+36-1 all
  Increase selected region in Emacs by semantic units

elpa-eyebrowse/stable,stable 0.7.8-2 all
  simple-minded way of managing window configs in Emacs

elpa-f/stable,stable 0.20.0-3 all
  modern API for working with files and directories in Emacs Lisp

elpa-faceup/stable,stable 0.0.4-5 all
  Regression test system for font-lock

elpa-fill-column-indicator/stable,stable 1.90-2.1 all
  graphically indicate the fill column

elpa-find-file-in-project/stable,stable 6.0.1-1 all
  quick access to project files in Emacs

elpa-flx/stable,stable 0.6.1-5 all
  sorting algorithm for fuzzy matching in Emacs

elpa-flx-ido/stable,stable 0.6.1-5 all
  allows Emacs Ido to use the flx sorting algorithm

elpa-flycheck/stable,stable 32~git.20200527.9c435db3-2 all
  modern on-the-fly syntax checking for Emacs

elpa-flycheck-package/stable,stable 0.13-1 all
  flycheck checker for Elisp package authors

elpa-folding/stable,stable 0+20200825.748-1 all
  folding-editor minor mode for Emacs

elpa-fountain-mode/stable,stable 2.8.5-1 all
  Emacs major mode for screenwriting in Fountain markup

elpa-fricas/stable,stable 1.3.6-6 all
  General purpose computer algebra system: emacs support

elpa-fsm/stable,stable 0.2.1-4 all
  state machine library

elpa-geiser/stable,stable 0.10-1 all
  enhanced Scheme interaction mode for Emacs

elpa-ggtags/stable,stable 0.8.13-2 all
  improved Emacs interface to GNU GLOBAL

elpa-git-annex/stable,stable 1.1-4 all
  Emacs-Integration für git-annex

elpa-git-auto-commit-mode/stable,stable 4.7.0-2 all
  Emacs Minor mode to automatically commit and push with git

elpa-git-commit/stable,stable 2.99.0.git0957.ge8c7bd03-1 all
  Major mode for editing git commit message

elpa-git-timemachine/stable,stable 4.11-1 all
  walk through git revisions of a file

elpa-gitlab-ci-mode/stable,stable 20190824.12.2-2 all
  Emacs mode for editing GitLab CI files

elpa-gnuplot-mode/stable,stable 1:0.7.0-2014-12-31-2 all
  Gnuplot mode for Emacs

elpa-go-mode/stable,stable 3:1.5.0-4 all
  Emacs mode for editing Go code

elpa-golden-ratio/stable,stable 1.0-6 all
  automatic resizing of Emacs windows to the golden ratio

elpa-goo/stable,stable 0.155+ds-4 all
  generic object-orientator (Emacs support)

elpa-goto-chg/stable,stable 1.7.3-1 all
  navigate the point to the most recent edit in the buffer

elpa-graphviz-dot-mode/stable,stable 0.4.2-2 all
  Emacs mode for the dot-language used by graphviz.

elpa-haskell-mode/stable,stable 17.2-3 all
  Emacs-Hauptmodus zum Editieren von Haskell

elpa-helm/stable,stable 3.7.0-2 all
  Emacs incremental completion and selection narrowing framework

elpa-helm-ag/stable,stable 0.59-1 all
  Silver Searcher integration with Emacs Helm

elpa-helm-core/stable,stable 3.7.0-2 all
  Emacs Helm library files

elpa-helm-org/stable,stable 1.0-2 all
  Emacs Helm for Org-mode headlines and keywords completion

elpa-highlight-indentation/stable,stable 0.7.0-5 all
  highlight the indentation level in Emacs buffers

elpa-ht/stable,stable 2.3-1 all
  hash table library for Emacs

elpa-htmlize/stable,stable 1.55-1 all
  convert buffer text and decorations to HTML

elpa-hydra/stable,stable 0.15.0-3 all
  make Emacs bindings that stick around

elpa-ibuffer-projectile/stable,stable 0.3-1 all
  group buffers in ibuffer list by Projectile project

elpa-ibuffer-vc/stable,stable 0.11-1 all
  group ibuffer list by VC project and show VC status

elpa-ido-completing-read+/stable,stable 4.13-2 all
  completing-read-function using ido

elpa-iedit/stable,stable 0.9.9.9-5 all
  edit multiple regions in the same way simultaneously

elpa-imenu-list/stable,stable 0.8-3 all
  show the current Emacs buffer's imenu entries in a separate window

elpa-initsplit/stable,stable 1.8+3+gc941d43-3 all
  code to split customizations into different files

elpa-irony/stable,stable 1.4.0+7.g76fd37f-1 all
  Emacs C/C++ minor mode powered by libclang

elpa-ivy/stable,stable 0.13.0-1 all
  generic completion mechanism for Emacs

elpa-ivy-hydra/stable,stable 0.13.0-1 all
  additional key bindings for Emacs Ivy

elpa-jabber/stable,stable 0.8.92+git98dc8e-6 all
  Jabber client for Emacsen

elpa-jedi/stable,stable 0.2.8-1 all
  Python auto-completion for Emacs

elpa-jedi-core/stable,stable 0.2.8-1 all
  common code of jedi.el and company-jedi.el

elpa-jinja2-mode/stable,stable 0.2+git20200624.159558e-1 all
  Emacs major mode for editing jinja2 code

elpa-js2-mode/stable,stable 0~20201220-1 all
  Emacs mode for editing Javascript programs

elpa-key-chord/stable,stable 0.6-5 all
  map pairs of simultaneously pressed keys to commands

elpa-kv/stable,stable 0.0.19+git20140108.7211484-4 all
  key/value data structure functions for Emacs Lisp

elpa-lbdb/stable,stable 0.49 all
  Little Brother's DataBase Emacs extensions

elpa-ledger/stable,stable 3.1.2~pre3+g5067e408-2 all
  command-line double-entry accounting program (emacs interface)

elpa-let-alist/stable,stable 1.0.6-2 all
  let-bind values of an assoc-list by their names in Emacs Lisp

elpa-linum-relative/stable,stable 0.6-2.1 all
  display relative line number in Emacs

elpa-load-relative/stable,stable 1.3.1-3 all
  relative file load (within a multi-file Emacs package)

elpa-loop/stable,stable 1.3-2.1 all
  friendly imperative loop structures for Emacs Lisp

elpa-lsp-haskell/stable,stable 1.0.20201011-1 all
  Haskell support for lsp-mode

elpa-lsp-mode/stable,stable 7.0.1-2 all
  Emacs client/library for the Language Server Protocol

elpa-lua-mode/stable,stable 20201010-1 all
  Emacs major-mode for editing Lua programs

elpa-macaulay2/stable,stable 1.17.1+ds-2 all
  Software system for algebraic geometry research (Emacs package)

elpa-magit/stable,stable 2.99.0.git0957.ge8c7bd03-1 all
  Emacs-Schnittstelle für Git

elpa-magit-forge/stable,stable 0.1.0+git20200714.639ce51-3 all
  Work with Git forges from the comfort of Magit

elpa-magit-popup/stable,stable 2.13.2-1 all
  Use popup like Magit

elpa-mailscripts/stable,stable 0.23-1 all
  Emacs functions for accessing tools in the mailscripts package

elpa-makey/stable,stable 0.3-4 all
  flexible context menu system

elpa-markdown-mode/stable,stable 2.4-1 all
  mode for editing Markdown-formatted text files in GNU Emacs

elpa-markdown-toc/stable,stable 0.1.5-1 all
  Emacs TOC (table of contents) generator for markdown files

elpa-meson-mode/stable,stable 0.3-1 all
  Major mode for the Meson build system files

elpa-message-templ/stable,stable 0.3.20161104-3 all
  templates for Emacs message-mode

elpa-migemo/stable,stable 1.9.2-3 all
  Japanese incremental search with Romaji on Emacsen

elpa-mocker/stable,stable 0.5.0-1 all
  mocking framework for Emacs

elpa-modus-themes/stable,stable 1.0.2-1 all
  set of accessible themes conforming with WCAG AAA accessibility standard

elpa-monokai-theme/stable,stable 3.5.3-3 all
  fruity color theme for Emacs

elpa-muse/stable,stable 3.20+dfsg-6 all
  author and publish projects using Wiki-like markup

elpa-mutt-alias/stable,stable 1.5-4 all
  Emacs package to lookup and insert expanded Mutt mail aliases

elpa-muttrc-mode/stable,stable 1.2.1-3 all
  Emacs major mode for editing muttrc

elpa-neotree/stable,stable 0.5.2-3 all
  directory tree sidebar for Emacs that is like NERDTree for Vim

elpa-no-littering/stable,stable 1.2.1-1 all
  help keeping ~/.emacs.d clean

elpa-noflet/stable,stable 0.0.15-5 all
  Emacs Lisp noflet macro for dynamic, local advice

elpa-nose/stable,stable 0.1.1-5 all
  easy Python test running in Emacs

elpa-notmuch/stable,stable 0.31.4-2 all
  thread-based email index, search and tagging (emacs interface)

elpa-nov/stable,stable 0.3.0-1 all
  featureful EPUB (ebook) reader mode for Emacs

elpa-olivetti/stable,stable 1.11.3-1 all
  Emacs minor mode to more comfortably read and write long-lined prose

elpa-openwith/stable,stable 0.8g-5 all
  seamlessly open files in external programs with Emacs

elpa-org/stable,stable 9.4.0+dfsg-1 all
  Keep notes, maintain ToDo lists, and do project planning in emacs

elpa-org-drill/stable,stable 2.7.0+20200412+dfsg1-2 all
  emacs org-mode contrib for self-testing using spaced repetition

elpa-org-roam/stable,stable 1.2.3-2 all
  non-hierarchical note-taking with Emacs Org-mode

elpa-orgalist/stable,stable 1.12-2 all
  Manage Org-like lists in non-Org Emacs buffers

elpa-package-lint/stable,stable 0.13-1 all
  linting library for Elisp package authors

elpa-package-lint-flymake/stable,stable 0.13-1 all
  package-lint Flymake backend

elpa-paredit/stable,stable 24-5 all
  Emacs minor mode for structurally editing Lisp code

elpa-paredit-everywhere/stable,stable 0.4-4 all
  cut-down version of paredit for non-lisp buffers

elpa-parent-mode/stable,stable 2.3-5 all
  get major mode's parent modes

elpa-parsebib/stable,stable 2.3.1-4 all
  Emacs Lisp library for parsing .bib files

elpa-pcre2el/stable,stable 1.8-4 all
  Emacs mode to convert between PCRE, Emacs and rx regexp syntax

elpa-pdf-tools/stable,stable 1.0~20200512-2 all
  Display and interact with pdf in Emacs

elpa-pdf-tools-server/stable 1.0~20200512-2 amd64
  server for Emacs's pdf-tools

elpa-persist/stable,stable 0.4+dfsg-2 all
  persist variables between Emacs Sessions

elpa-perspective/stable,stable 2.2-3 all
  tagged workspaces in Emacs

elpa-pg/stable,stable 0.13+git.20130731.456516ec-2 all
  Emacs Lisp interface for PostgreSQL

elpa-php-mode/stable,stable 1.23.0-1 all
  PHP Mode for GNU Emacs

elpa-pip-requirements/stable,stable 0.5-3 all
  major mode for editing pip requirements files

elpa-pkg-info/stable,stable 0.6-6 all
  provide information about Emacs packages

elpa-pod-mode/stable,stable 1.03-3 all
  Emacs major mode for editing .pod files

elpa-popup/stable,stable 0.5.8-1 all
  visual popup user interface library for Emacs

elpa-pos-tip/stable,stable 0.4.6+git20191227-2 all
  Show tooltip at point

elpa-powerline/stable,stable 2.4-4 all
  Emacs version of the Vim powerline

elpa-project/stable,stable 0.5.2-2 all
  Emacs library for operations on the current project

elpa-projectile/stable,stable 2.1.0-1 all
  project interaction library for Emacs

elpa-protobuf-mode/stable,stable 3.12.4-1 all
  Emacs-Addon zur Bearbeitung von Protocol Buffers

elpa-ps-ccrypt/stable,stable 1.11-2 all
  Emacs addon for working with files encrypted with ccrypt

elpa-puppet-mode/stable,stable 0.4-2 all
  Emacs major mode for Puppet manifests

elpa-python-environment/stable,stable 0.0.2-6 all
  virtualenv API for Emacs Lisp

elpa-pyvenv/stable,stable 1.21+git20201124.37e7cb1-1 all
  Python virtual environment interface

elpa-qml-mode/stable,stable 0.4-4 all
  Emacs major mode for editing QT Declarative (QML) code

elpa-queue/stable,stable 0.2-3 all
  queue data structure for Emacs Lisp

elpa-racket-mode/stable,stable 20201227git0-3 all
  emacs support for editing and running racket code

elpa-rainbow-delimiters/stable,stable 2.1.3-5 all
  Emacs mode to colour-code delimiters according to their depth

elpa-redtick/stable,stable 00.01.02+git20170220.e6d2e9b+dfsg-4 all
  tiny pomodoro timer for Emacs

elpa-relint/stable,stable 1.19-1 all
  Emacs Lisp regexp mistake finder

elpa-restart-emacs/stable,stable 0.1.1-4 all
  Emacs aus Emacs heraus neu starten

elpa-rich-minority/stable,stable 1.0.3-2 all
  clean-up and beautify the list of minor-modes in Emacs' mode-line

elpa-rtags/stable,stable 2.38-3 all
  emacs front-end for RTags

elpa-rust-mode/stable,stable 0.4.0-2 all
  Major Emacs mode for editing Rust source code

elpa-s/stable,stable,now 1.12.0-4 all  [Installiert,automatisch]
  string manipulation library for Emacs

elpa-scala-mode/stable,stable 1:1.1.0-2 all
  Emacs major mode for editing scala source code

elpa-seq/stable,stable 2.22-1 all
  sequence manipulation functions for Emacs Lisp

elpa-sesman/stable,stable 0.3.4-2 all
  session manager for Emacs IDEs

elpa-session/stable,stable 2.4b-3 all
  use variables, registers and buffer places across sessions

elpa-shut-up/stable,stable 0.3.3-1 all
  Emacs Lisp macros to quieten Emacs

elpa-smart-mode-line/stable,stable 2.13-2 all
  powerful and beautiful mode-line for Emacs

elpa-smart-mode-line-powerline-theme/stable,stable 2.13-2 all
  Smart Mode Line themes that use Emacs Powerline

elpa-smeargle/stable,stable 0.03-5 all
  highlight region by last updated time

elpa-smex/stable,stable 3.0-6 all
  M-x interface for Emacs with Ido-style fuzzy matching

elpa-sml-mode/stable,stable 6.10-1 all
  Emacs major mode for editing Standard ML programs

elpa-solarized-theme/stable,stable 1.3.1-1 all
  port of Solarized theme to Emacs

elpa-spinner/stable,stable 1.7.3-3 all
  spinner for the Emacs modeline for operations in progress

elpa-suggest/stable,stable 0.7-3 all
  discover Emacs Lisp functions based on examples

elpa-super-save/stable,stable 0.3.0-3 all
  auto-save buffers, based on your activity

elpa-swiper/stable,stable 0.13.0-1 all
  alternative to Emacs' isearch, with an overview

elpa-sxiv/stable,stable 0.3.3-1 all
  run the sxiv image viewer

elpa-systemd/stable,stable 1.6-2.1 all
  major mode for editing systemd units

elpa-tabbar/stable,stable 2.2-4 all
  Emacs minor mode that displays a tab bar at the top

elpa-transient/stable,stable 0.2.0.30.g4d44d08-2 all
  Emacs key and popup interface for complex keybindings

elpa-transient-doc/stable,stable 0.2.0.30.g4d44d08-2 all
  Emacs key and popup interface for complex keybindings

elpa-transmission/stable,stable 0.12.2-1 all
  Emacs interface to a Transmission session

elpa-tuareg/stable,stable 1:2.2.0-1 all
  emacs-mode for OCaml programs

elpa-undercover/stable,stable 0.8.0-1 all
  test coverage library for Emacs Lisp

elpa-undo-tree/stable,stable 0.7.4-1 all
  Emacs minor mode for handling undo history as tree

elpa-use-package/stable,stable 2.4.1-1 all
  configuration macro for simplifying your .emacs

elpa-uuid/stable,stable 0.0.3~git20120910.1519bfe-3 all
  UUID/GUID library for Emacs Lisp

elpa-vala-mode/stable,stable 0.1-8 all
  Emacs editor major mode for vala source code

elpa-vc-fossil/stable,stable 2020.09.20-4 all
  Emacs VC backend for the Fossil Version Control system

elpa-verbiste/stable,stable 0.1.47-1 all
  Konjugationssystem für Französisch und Italienisch - Emacs-Erweiterung

elpa-vimish-fold/stable,stable 0.2.3-5 all
  fold text in GNU Emacs like in Vim

elpa-virtualenvwrapper/stable,stable 0.2.0-2 all
  featureful virtualenv tool for Emacs

elpa-visual-fill-column/stable,stable 2.3-1 all
  Emacs mode that wraps visual-line-mode buffers at fill-column

elpa-visual-regexp/stable,stable 1.1.2-2 all
  in-buffer visual feedback while using Emacs regexps

elpa-volume/stable,stable 1.0+git.20201002.afb75a5-3 all
  tweak your sound card volume from Emacs

elpa-wc-mode/stable,stable 1.4-1 all
  display a word count in the Emacs modeline

elpa-web-mode/stable,stable 17.0.2-1 all
  major emacs mode for editing web templates

elpa-websocket/stable,stable 1.13-1 all
  Emacs WebSocket client and server

elpa-weechat/stable,stable 0.5.0-5 all
  Chat via WeeChat's relay protocol in Emacs.

elpa-wgrep/stable,stable 2.3.2+9.gf0ef9bf-2 all
  edit multiple Emacs buffers using a master grep pattern buffer

elpa-wgrep-ack/stable,stable 2.3.2+9.gf0ef9bf-2 all
  edit multiple Emacs buffers using a master ack pattern buffer

elpa-wgrep-ag/stable,stable 2.3.2+9.gf0ef9bf-2 all
  edit multiple Emacs buffers using a master ag pattern buffer

elpa-wgrep-helm/stable,stable 2.3.2+9.gf0ef9bf-2 all
  edit multiple Emacs buffers with a helm-grep-mode buffer

elpa-which-key/stable,stable 3.5.1-1 all
  display available keybindings in popup

elpa-with-editor/stable,stable 3.0.2-1 all
  call program using Emacs as $EDITOR

elpa-with-simulated-input/stable,stable 2.4+git20200216.29173588-1 all
  macro to simulate user input non-interactively

elpa-world-time-mode/stable,stable 0.0.6-4 all
  Emacs mode to compare timezones throughout the day

elpa-writegood-mode/stable,stable 2.0.3-3 all
  Emacs minor mode that provides hints for common English writing problems

elpa-ws-butler/stable,stable 0.6-4 all
  unobtrusively remove trailing whitespace in Emacs

elpa-xcite/stable,stable 1.60-7 all
  exciting cite utility for Emacsen

elpa-xcscope/stable,stable 1.5-1.1 all
  Interactively examine a C program source in emacs

elpa-xml-rpc/stable,stable 1.6.12-4 all
  Emacs Lisp XML-RPC client

elpa-xr/stable,stable 1.21-1 all
  convert string regexp to rx notation

elpa-xref/stable,stable 1.0.2-2 all
  Library for cross-referencing commands in Emacs

elpa-yaml-mode/stable,stable 0.0.15-1 all
  Emacs major mode for YAML files

elpa-yasnippet/stable,stable 0.14.0+git20200603.5cbdbf0d-1 all
  template system for Emacs

elpa-yasnippet-snippets/stable,stable 0.23-1 all
  Andrea Crotti's official YASnippet snippets

elpa-zenburn-theme/stable,stable 2.6-3 all
  low contrast color theme for Emacs

elpa-ztree/stable,stable 1.0.5-4 all
  text mode directory tree

elscreen/stable,stable 1.4.6-5.3 all
  Screen für Emacse

emacs/stable-security,stable-security,now 1:27.1+1-3.1+deb11u2 all  [installiert]
  Editor GNU Emacs (Metapaket)

emacs-bin-common/stable-security,now 1:27.1+1-3.1+deb11u2 amd64  [Installiert,automatisch]
  Editor GNU Emacs - gemeinsame, architekturabhängige Dateien

emacs-calfw/stable,stable 1.6+git20180118-1.1 all
  calendar framework for Emacs

emacs-calfw-howm/stable,stable 1.6+git20180118-1.1 all
  calendar framework for Emacs (howm add-on)

emacs-common/stable-security,stable-security,now 1:27.1+1-3.1+deb11u2 all  [Installiert,automatisch]
  Editor GNU Emacs - gemeinsame, architekturunabhängige Infrastruktur

emacs-common-non-dfsg/stable,stable 1:27.1+1-2 all
  GNU Emacs common non-DFSG items, including the core documentation

emacs-el/stable-security,stable-security,now 1:27.1+1-3.1+deb11u2 all  [Installiert,automatisch]
  LISP-Dateien (.el) für den Editor GNU Emacs

emacs-goodies-el/stable,stable 42.3 all
  Miscellaneous add-ons for Emacs

emacs-gtk/stable-security,now 1:27.1+1-3.1+deb11u2 amd64  [Installiert,automatisch]
  Editor GNU Emacs (mit Unterstützung für eine GTK+-Benutzeroberfläche)

emacs-intl-fonts/stable,stable 1.2.1-10.1 all
  fonts to allow multilingual PostScript printing from Emacs

emacs-jabber/stable,stable 0.8.92+git98dc8e-6 all
  Transition package, emacs-jabber to elpa-jabber

emacs-lucid/stable-security 1:27.1+1-3.1+deb11u2 amd64
  GNU Emacs editor (with Lucid GUI support)

emacs-mozc/stable 2.26.4220.100+dfsg-4 amd64
  Mozc for Emacs

emacs-mozc-bin/stable 2.26.4220.100+dfsg-4 amd64
  Helper module for emacs-mozc

emacs-nox/stable-security 1:27.1+1-3.1+deb11u2 amd64
  Editor GNU Emacs (ohne Unterstützung einer grafischen Oberfläche)

emacs-window-layout/stable,stable 1.4-2.1 all
  window layout manager for emacs

emacsen-common/stable,stable,now 3.0.4 all  [Installiert,automatisch]
  Gemeinsame Funktionen aller Emacs-Varianten

emacspeak/stable,stable 53.0+dfsg-1 all
  Sprachausgabe-Schnittstelle für Emacs

emacspeak-espeak-server/stable 53.0+dfsg-1 amd64
  espeak synthesis server for emacspeak

emacspeak-ss/stable 1.12.1-8 amd64
  Emacspeak-Sprachserver für verschiedene Synthesizer

emms/stable 5.1-1+b1 amd64
  Emacs-Multimediasystem

erlang/stable,stable 1:23.2.6+dfsg-1+deb11u1 all
  Simultane, verteilte und funktionelle Echtzeitsprache

erlang-mode/stable,stable 1:23.2.6+dfsg-1+deb11u1 all
  Emacs-Haupteditiermodus für Erlang

erlang-tools/stable 1:23.2.6+dfsg-1+deb11u1 amd64
  Verschiedene Werkzeuge für Erlang/OTP

ess/stable,stable 18.10.2-2 all
  Übergangspaket für den Wechsel von ess zu elpa-ess

etktab/stable,stable 3.2-13 all
  ASCII Gitarrengriffschrift-Editor

eweouz/stable 0.12+b1 amd64
  Emacs interface to Evolution Data Server

exuberant-ctags/stable 1:5.9~svn20110310-14 amd64
  Erzeugt Indexdateien von Quelltextdefinitionen

fetchmail/stable 6.4.16-4+deb11u1 amd64
  SSL-fähiger E-Mail-Sammler/-Versender für POP3, APOP und IMAP

findent/stable 3.1.7-1 amd64
  indents/converts Fortran sources

flim/stable,stable 1:1.14.9+0.20201117-2 all
  library about internet message for emacsen

flycheck-doc/stable,stable 32~git.20200527.9c435db3-2 all
  modern on-the-fly syntax checking for Emacs - documentation

fortran-language-server/stable,stable 1.12.0-1 all
  Fortran Language Server for the Language Server Protocol

geiser/stable,stable 0.10-1 all
  Transition Package, geiser to elpa-geiser

gettext-el/stable,stable 0.21-4 all
  Emacs-Modus zur Bearbeitung der gettext-PO-Dateien

gir1.2-kkc-1.0/stable 0.3.5-7 amd64
  GObject introspection data for libkkc

git-el/stable,stable,stable-security,stable-security 1:2.30.2-1+deb11u2 all
  Schnelles, skalierbares, verteiltes Versionskontrollsystem (Emacs-Unterstützung)

global/stable 6.6.5-1 amd64
  Werkzeuge zum Suchen und Browsen in Quelltext

gmult/stable 8.0-2+b1 amd64
  Finden Sie heraus, welche Buchstaben welche Ziffern darstellen!

gnu-smalltalk-el/stable,stable 3.2.5-1.3 all
  GNU Smalltalk Emacs front-end

gnuplot-mode/stable,stable 1:0.7.0-2014-12-31-2 all
  Transition Package, gnuplot-mode to elpa-gnuplot-mode

gnuserv/stable 3.12.8-7+b2 amd64
  Erlaubt das Anbinden an einen bereits laufenden Emacs

goby/stable,stable 1.1+0.20180214-5 all
  WYSIWYG presentation tool for Emacs

golang-mode/stable,stable 3:1.5.0-4 all
  Emacs mode for editing Go code -- transitional package

gramadoir/stable,stable 0.7-4.1 all
  Irish language grammar checker (integration scripts)

gri-el/stable,stable 2.12.27-1.1~deb11u1 all
  Emacs major-mode for gri, a language for scientific graphics

haml-elisp/stable,stable 1:3.1.0-3.2 all
  Emacs Lisp mode for the Haml markup language

howm/stable,stable 1.4.7-1 all
  Note-taking tool on Emacs

id-utils/stable 4.6.28-20200521ss15dab+b1 amd64
  Schnelles Werkzeug mit hohem Durchsatz für Bezeichnerdatenbanken

idl-font-lock-el/stable,stable 1.5-9.1 all
  OMG IDL Schriften-Sperrung für Emacs

idn/stable 1.33-3 amd64
  Befehlszeilen- und Emacs-Schnittstelle für GNU Libidn

ilisp/stable,stable 5.12.0+cvs.2004.12.26-28 all
  Emacs-Schnittstelle zu LISP-Implementationen

ilisp-doc/stable,stable 5.12.0+cvs.2004.12.26-28 all
  Dokumentation für das Paket ILISP

info2man/stable,stable 1.1-10 all
  Konvertiert GNU-Info-Dateien in POD oder Handbuchseiten

initz/stable,stable 0.0.11+20030603cvs-17.3 all
  Nutzung von verschieden Initdateien für Emacsen

inotify-hookable/stable,stable 0.09-2 all
  blocking command-line interface to inotify

irony-server/stable 1.4.0+7.g76fd37f-1 amd64
  Emacs C/C++ minor mode powered by libclang (server)

ispell/stable,now 3.4.02-2 amd64  [Installiert,automatisch]
  Internationales Ispell (eine interaktive Schreibkorrektur)

jed/stable 1:0.99.19-8 amd64
  Editor für Programmierer (Textmodus-Version)

jedit/stable,stable 5.5.0+dfsg-2 all
  Plugin-based editor for programmers

joe/stable 4.6-1+b1 amd64
  benutzerfreundlicher Vollbild-Texteditor

jove/stable 4.17.3.6-2 amd64
  Jonathans Version von Emacs - ein kompakter, mächtiger Editor

js2-mode/stable,stable 0~20201220-1 all
  Emacs mode for editing Javascript programs (dummy package)

jupp/stable 3.1.40-1 amd64
  user friendly full screen text editor

kdesdk-scripts/stable,stable 4:20.12.0-1 all
  Skripte und Datendateien für die Entwicklung

latex-cjk-common/stable 4.8.4+git20170127-3 amd64
  LaTeX-Makropaket für CJK (Chinesisch/Japanisch/Koreanisch)

ledit/stable,stable 2.04-5 all
  Zeileneditor für interaktive Programme

libconfig-find-perl/stable,stable 0.31-1.1 all
  module to search configuration files using OS dependent heuristics

libghc-agda-dev/stable 2.6.1-1+b2 amd64
  Abhängig typisierte, funktionale Programmiersprache

libghc-agda-doc/stable,stable 2.6.1-1 all
  Abhängig typisierte, funktionale Programmiersprache - Dokumentation

libghc-pandoc-dev/stable 2.9.2.1-1+b1 amd64
  general markup converter - libraries

libghc-pandoc-doc/stable,stable 2.9.2.1-1 all
  general markup converter - library documentation

libghc-pandoc-prof/stable 2.9.2.1-1+b1 amd64
  general markup converter - profiling libraries

libghc-yi-keymap-emacs-dev/stable 0.19.0-1 amd64
  Emacs keymap for Yi editor

libghc-yi-keymap-emacs-doc/stable,stable 0.19.0-1 all
  Emacs keymap for Yi editor; documentation

libghc-yi-keymap-emacs-prof/stable 0.19.0-1 amd64
  Emacs keymap for Yi editor; profiling libraries

libjline2-java/stable,stable 2.14.6-4 all
  console input handling in Java

libkkc-common/stable,stable 0.3.5-7 all
  Japanese Kana Kanji input library - common data

libkkc-data/stable 0.2.7-4 amd64
  language model data for libkkc

libkkc-dev/stable 0.3.5-7 amd64
  Japanese Kana Kanji input library - development files

libkkc-utils/stable 0.3.5-7 amd64
  Japanese Kana Kanji input library - testing utility

libkkc2/stable 0.3.5-7 amd64
  Japanese Kana Kanji input library

liblatex-table-perl/stable,stable 1.0.6-3.1 all
  Perl extension for the automatic generation of LaTeX tables

libledit-ocaml-dev/stable 2.04-5 amd64
  OCaml line editor library

libocp-indent-ocaml/stable 1.8.2-1+b1 amd64
  OCaml indentation tool for emacs and vim - libraries

libocp-indent-ocaml-dev/stable 1.8.2-1+b1 amd64
  OCaml indentation tool for emacs and vim - development libraries

libparse-exuberantctags-perl/stable 1.02-1+b8 amd64
  exuberant ctags parser for Perl

libpcre-ocaml/stable 7.4.6-1+b1 amd64
  OCaml bindings for PCRE (runtime)

libpcre-ocaml-dev/stable 7.4.6-1+b1 amd64
  OCaml bindings for PCRE (Perl Compatible Regular Expression)

libproc-invokeeditor-perl/stable,stable 1.13-1.1 all
  Perl extension for starting a text editor

libre-ocaml-dev/stable 1.9.0-1+b1 amd64
  regular expression library for OCaml

librep-dev/stable 0.92.5-3+b6 amd64
  Entwicklungsbibliotheken und Header für librep

librep16/stable 0.92.5-3+b6 amd64
  embedded lisp command interpreter library

librobert-hooke-clojure/stable,stable 1.3.0-4 all
  Function wrapper library for Clojure

libtext-findindent-perl/stable,stable 0.11-1 all
  module to heuristically determine indentation style

libutop-ocaml/stable 2.7.0-2 amd64
  improved OCaml toplevel (runtime library)

libutop-ocaml-dev/stable 2.7.0-2 amd64
  improved OCaml toplevel (development tools)

liece/stable,stable 2.0+0.20030527cvs-12 all
  IRC (Internet Relay Chat) client for Emacs

liece-dcc/stable 2.0+0.20030527cvs-12+b1 amd64
  DCC-Programm für liece

liquidsoap-mode/stable,stable 1.4.3-3 all
  Emacs mode for editing Liquidsoap code

lisaac-mode/stable,stable 1:0.39~rc1-3.1 all
  Emacs mode for editing Lisaac programs

lookup-el/stable,stable 1.4.1-20 all
  emacsen interface to electronic dictionaries

lsdb/stable,stable 0.11-10.2 all
  die Lovely Sister Database (email rolodex) für Emacs

lyskom-elisp-client/stable,stable 0.48+git.20200923.ec349ff4-3 all
  Emacs-CLient für LysKOM

malaga-mode/stable,stable 7.12-7.1 all
  System für automatische Sprachanalyse - Emacs-Modus

maxima-emacs/stable,stable 5.44.0-3 all
  Computeralgebrasystem -- Emacs-Schnittstelle

mew/stable,stable 1:6.8-13 all
  Mailreader mit PGP/MIME-Unterstützung für Emacs

mew-beta/stable,stable 7.0.50~6.8+0.20210131-2 all
  mail reader supporting PGP/MIME for Emacs (development version)

mg/stable,now 20200723-1 amd64  [installiert]
  Mikroskopischer Editor im Stil von GNU Emacs

mgp/stable 1.13a+upstream20090219-12 amd64
  MagicPoint — ein X11-basiertes Präsentationsprogramm

mhc/stable,stable 1.2.4-2 all
  schedule management tool for Emacs

mhc-utils/stable,stable 1.2.4-2 all
  utilities for the MHC schedule management system

midge/stable,stable 0.2.41-2.1 all
  Ein Text-zu-MIDI-Programm

minlog/stable,stable 4.0.99.20100221-7 all
  Proof assistant based on first order natural deduction calculus

mit-scheme/stable 10.1.11-2 amd64
  MIT/GNU Scheme development environment

mit-scheme-dbg/stable 10.1.11-2 amd64
  MIT/GNU Scheme debugging files

mit-scheme-doc/stable,stable 10.1.11-2 all
  MIT/GNU Scheme documentation

mksh/stable 59c-9+b2 amd64
  MirBSD Korn Shell

mmm-mode/stable,stable 0.5.8-1 all
  Mehrfacher »Major Mode« für Emacs

mpqc-support/stable 2.3.1-21 amd64
  Massiv-paralleles Quantenchemieprogramm (Hilfsprogramme)

mu-cite/stable,stable 8.1+0.20201103-2 all
  message citation utility for emacsen

mu4e/stable,stable 1.4.15-1 all
  E-Mail-Client für Emacs, basiert auf mu (maildir-utils)

nescc/stable 1.3.5-1.1 amd64
  Programming Language for Deeply Networked Systems

ng-cjk/stable 1.5~beta1-9 amd64
  Nihongo MicroGnuEmacs with CJK support

ng-cjk-canna/stable 1.5~beta1-9 amd64
  Nihongo MicroGnuEmacs with CJK and Canna support

ng-common/stable,stable 1.5~beta1-9 all
  Common files used by ng-* packages

ng-latin/stable 1.5~beta1-9 amd64
  Nihongo MicroGnuEmacs with Latin support

nmh/stable 1.7.1-7 amd64
  Programme für die Verarbeitung von E-Mails

nomarch/stable 1.4-4 amd64
  Entpackt .ARC- und .ARK-MS-DOS-Archive

notmuch-addrlookup/stable 9-2 amd64
  Address lookup tool for Notmuch

ocaml-core/stable,stable 4.08.1.2 all
  OCaml core tools (metapackage)

ocp-indent/stable 1.8.2-1+b1 amd64
  OCaml indentation tool for emacs and vim - runtime

oneliner-el/stable,stable 0.3.6-9.1 all
  extensions of Emacs standard shell-mode

org-mode/stable,stable 9.4.0+dfsg-1 all
  Transition Package, org-mode to elpa-org

org-mode-doc/stable,stable 9.4.0-2 all
  keep notes, maintain ToDo lists, and do project planning in emacs

org-roam-doc/stable,stable 1.2.3-2 all
  non-hierarchical note-taking with Emacs Org-mode -- documentation

otags/stable 4.05.1-2+b2 amd64
  tags file generator for OCaml

pandoc/stable 2.9.2.1-1+b1 amd64
  general markup converter

pandoc-data/stable,stable 2.9.2.1-1 all
  general markup converter - data files

post-el/stable,stable 1:2.6-2 all
  Emacs-Mode zur Mailbearbeitung

projectile-doc/stable,stable 2.1.0-1 all
  project interaction library for Emacs - documentation

proofgeneral/stable,stable 4.4.1~pre170114-1.2 all
  Generisches Frontend für Beweisassistenten

proofgeneral-doc/stable,stable 4.4.1~pre170114-1.2 all
  Generisches Frontend für Beweisassistenten - Dokumentation

psgml/stable,stable 1.4.0-12 all
  Emacs major mode for editing SGML documents

pylint/stable,stable 2.7.2-3 all
  Statisches Prüfprogramm für Python-3-Code und Generator von UML-Diagrammen

pymacs/stable,stable 0.25-3 all
  Schnittstelle zwischen Emacs Lisp und Python

python3-editor/stable,stable 1.0.3-2 all
  programmatically open an editor, capture the result - Python 3.x

python3-epc/stable,stable 0.0.5-3 all
  RPC stack for Emacs Lisp (Python3 version)

python3-readlike/stable,stable 0.1.3-1.1 all
  GNU Readline-like line editing module

quilt-el/stable,stable 0.66-2.1 all
  Einfache Emacs-Schnittstelle zu quilt

r-cran-progress/stable,stable 1.2.2-2 all
  GNU R terminal progress bars

rabbit-mode/stable,stable 3.0.0-4 all
  Emacs-lisp rabbit-mode for writing RD document using Rabbit

rail/stable,stable 1.2.14-2 all
  Replace Agent-string Internal Library

rdtool-elisp/stable,stable 0.6.38-4 all
  Emacs-lisp rd-Modus zum Erstellen von RD-Dokumenten

remembrance-agent/stable 2.12-7+b2 amd64
  Emacs-Modus zum Finden relevanter Texte

rep/stable 0.92.5-3+b6 amd64
  Befehlsinterpreter für Lisp

rep-doc/stable,stable 0.92.5-3 all
  Dokumentation für den Lisp-Befehlsinterpreter

reposurgeon/stable 4.25-1+b4 amd64
  Tool for editing version-control repository history

riece/stable,stable 9.0.0-11 all
  IRC client for Emacs

rtags/stable 2.38-3 amd64
  C/C++ client/server indexer with integration for Emacs

ruby-github-markup/stable,stable 1.7.0+dfsg-3 all
  GitHub Markup rendering

ruby-notiffany/stable,stable 0.1.3-1 all
  Wrapper libray for most popular notification libraries

ruby-org/stable,stable 0.9.12-2 all
  Emacs org-mode parser for Ruby

sass-elisp/stable,stable 3.0.15-4.4 all
  Emacs Lisp mode for the Sass markup language

sawfish/stable 1:1.11.90-1.2 amd64
  X11-Fenstermanager

search-ccsb/stable,stable 0.5-5 all
  Suchprogramm für BibTeX

search-citeseer/stable,stable 0.3-3 all
  BibTeX search tool

select-xface/stable,stable 0.15-12 all
  utility for selecting X-Face on emacsen

semi/stable,stable 1.14.7~0.20201115-2 all
  library to provide MIME feature for emacsen

sepia/stable,stable 0.992-7 all
  Simple Emacs-Perl InterAction

singular-ui-emacs/stable 1:4.1.1-p2+ds-4+b2 amd64
  Computer Algebra System for Polynomial Computations -- emacs user interface

sisu/stable,stable 7.2.0-1 all
  documents - structuring, publishing in multiple formats and search

slime/stable,stable 2:2.26.1+dfsg-2 all
  Superior Lisp Interaction Mode for Emacs (client)

speechd-el/stable,stable 2.9-2 all
  Emacs speech client using Speech Dispatcher

speechd-el-doc-cs/stable,stable 2.9-2 all
  speechd-el Dokumentation in Tschechisch

speechd-up/stable 0.5~20110719-10 amd64
  Interface between Speech Dispatcher and SpeakUp

stow/stable,stable 2.3.1-1 all
  Organizer for /usr/local software packages

stumpwm/stable,stable 2:1.0.0-1 all
  Kachelnder und tastaturgesteuerter Fenstermanager geschrieben in Common Lisp

supercollider-emacs/stable,stable 1:3.11.2+repack-1 all
  SuperCollider mode for Emacs

sylpheed/stable 3.7.0-8 amd64
  Schlanker E-Mail-Client mit GTK+

t-code/stable,stable 2:2.3.1-9 all
  Japanese direct input method environment for emacsen

t-code-common/stable,stable 2:2.3.1-9 all
  Japanese direct input method environment - common files

tdiary-mode/stable,stable 5.1.5-1 all
  Emacs-Modus für Editieren von tDiary

tiarra-conf-el/stable,stable 20100212+r39209-9 all
  edit mode for tiarra.conf

timidity-el/stable,stable 2.14.0-8 all
  Emacs-Frontend für TiMidity++

tkcon/stable,stable 2:2.7.3-1 all
  Erweiterte interaktive Konsole für die Entwicklung mit Tcl

tmux/stable,now 3.1c-1+deb11u1 amd64  [installiert]
  Terminal-Multiplexer

tpp/stable,stable 1.3.1-8 all
  text presentation program

tuareg-mode/stable,stable 1:2.2.0-1 all
  transitional package, tuareg-mode to elpa-tuareg

tweak/stable 3.02-6 amd64
  Efficient text-mode hex editor

twittering-mode/stable,stable 3.1.0-1.2 all
  Twitter client for Emacs

txt2regex/stable,stable 0.9-3 all
  Ein Suchmuster-Zauberer auf Basis von bash2

tzc/stable 2.6.15-5.4+b1 amd64
  Einfacher Zephyr-Client

uim-el/stable 1:1.8.8-9 amd64
  Universal Input Method - Emacs-Frontend

uim-latin/stable,stable 1:1.8.8-9 all
  Universal Input Method - Latin script input support metapackage

universal-ctags/stable 0+git20200824-1.1 amd64
  build tag file indexes of source code definitions

utop/stable 2.7.0-2 amd64
  improved OCaml toplevel

vile/stable 9.8u-2 amd64
  VI wie Emacs - arbeitet wie vi

vile-common/stable,stable 9.8u-2 all
  VI Like Emacs - Support-Dateien für vile/xvile

vile-filters/stable 9.8u-2 amd64
  VI Like Emacs - Hervorhebungsfilter für vile/xvile

vim-voom/stable,stable 5.3-8 all
  Vim two-pane outliner

vm/stable,stable 8.2.0b-7 all
  mail user agent for Emacs

w3m-el/stable,stable 1.4.632+0.20181112-9 all
  simple Emacs interface of w3m

w3m-el-snapshot/stable,stable 1.4.632+0.20210201.2305.54c3ccd-1 all
  simple Emacs interface of w3m (development version)

whizzytex/stable,stable 1.3.7-1 all
  WYSIWYG emacs environment for LaTeX

windows-el/stable,stable 2.55-1 all
  window manager for GNU Emacs

wl/stable,stable 2.15.9+0.20190205-7 all
  mail/news reader supporting IMAP for emacsen

wl-beta/stable,stable 2.15.9+0.20210131-2 all
  mail/news reader supporting IMAP for emacsen (development version)

wnn7egg/stable,stable 1.02-9 all
  Wnn-nana-tamago -- EGG Input Method with Wnn7 for Emacsen

wordwarvi/stable 1.0.4-1 amd64
  retro-styled side-scrolling shoot'em up arcade game

wordwarvi-sound/stable,stable 1.0.4-1 all
  retro-styled side-scrolling shoot'em up arcade game [Sound Files]

x-face-el/stable,stable 1.3.6.24-18 all
  utility for displaying X-Face on emacsen

xcite/stable,stable 1.60-7 all
  Transition Package, xcite to elpa-xcite

xcscope-el/stable,stable 1.5-1.1 all
  Transition Package, xcscope-el to elpa-xcscope

xemacs21/stable,stable 21.4.24-9 all
  highly customizable text editor metapackage

xemacs21-basesupport/stable,stable 2009.02.17.dfsg.2-5 all
  Editor und »Spülbecken« -- compilierte Elisp-Hilfsdateien

xemacs21-basesupport-el/stable,stable 2009.02.17.dfsg.2-5 all
  Editor und Abfluss -- Elisp-Unterstützung (Quell-Dateien)

xemacs21-bin/stable 21.4.24-9 amd64
  sehr flexibler Texteditor -- benötigte Binärdateien

xemacs21-mule/stable 21.4.24-9 amd64
  sehr flexibler Texteditor -- Mule Binärdatei

xemacs21-mule-canna-wnn/stable 21.4.24-9 amd64
  Sehr flexibler Texteditor -- Mule-Binärdatei mit Unterstützung für Canna und Wnn

xemacs21-mulesupport/stable,stable 2009.02.17.dfsg.2-5 all
  Editor und Küchenspüle -- Mule elisp-Unterstützungs-Dateien

xemacs21-mulesupport-el/stable,stable 2009.02.17.dfsg.2-5 all
  Editor und Abfluss -- Elisp-Unterstützung (Quell-Dateien)

xemacs21-nomule/stable 21.4.24-9 amd64
  Sehr flexibler Texteditor -- nicht-Mule Binärdatei

xemacs21-support/stable,stable 21.4.24-9 all
  Sehr flexibler Texteditor -- architekturunabhängige Unterstützungsdateien

xemacs21-supportel/stable,stable 21.4.24-9 all
  highly customizable text editor -- non-required library files

xfonts-kapl/stable,stable 4.22.1-10.1 all
  APL-Zeichensätze für A+-Entwicklung

xfonts-terminus-oblique/stable,stable 4.48-3 all
  Oblique version of the Terminus font

xfonts-thai-etl/stable,stable 1:1.2.7-5 all
  Thai etl fonts for X

xfonts-thai-poonlap/stable,stable 1:1.2.7-5 all
  Poonlap Veerathanabutr's bitmap fonts for X

xjed/stable 1:0.99.19-8 amd64
  Ein Editor für Programmierer (X11‐Version)

xstow/stable 1.0.2-1 amd64
  Extended replacement of GNU Stow

xvile/stable 9.8u-2 amd64
  VI Like Emacs - VI-artiger Editor (X11)

yasnippet/stable,stable 0.14.0+git20200603.5cbdbf0d-1 all
  transition Package, yasnippet to elpa-yasnippet

yasr/stable 0.6.9-10 amd64
  General-purpose console screen reader

yatex/stable,stable 1.82-1 all
  Ein weiterer TeX-Modus für Emacs

yc-el/stable 5.0.0-8.1 amd64
  Noch ein weiterer Canna-Client für Emacsen

yorick/stable 2.2.04+dfsg1-12 amd64
  interpreted language and scientific graphics

zile/stable 2.4.15-2 amd64
  Editor mit sehr kleiner Emacs-Untermenge

*/
func (pk *PkgApt) Search(query string) ([]Package, error) {
	const cmdSearch = cmdApt // "/usr/bin/apt-cache"
	var (
		err              error
		cmd              *exec.Cmd
		pipeOut, pipeErr io.ReadCloser
		bufOut, bufErr   bytes.Buffer
	)

	cmd = exec.Command(cmdSearch, "search", query)

	if pipeOut, err = cmd.StdoutPipe(); err != nil {
		pk.log.Printf("[ERROR] Cannot get stdout pipe from Cmd: %s\n",
			err.Error())
		return nil, err
	} else if pipeErr, err = cmd.StderrPipe(); err != nil {
		pk.log.Printf("[ERROR] Cannot get stderr pipe from Cmd: %s\n",
			err.Error())
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		pk.log.Printf("[ERROR] Error starting command: %s\n",
			err.Error())
		return nil, err
	}

	io.Copy(&bufOut, pipeOut) // nolint: errcheck
	io.Copy(&bufErr, pipeErr) // nolint: errcheck

	if err = cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			// FIXME Do something with stderr output!
			pk.log.Printf("[ERROR] Failed to wait for command: %s\n",
				err.Error())
			return nil, err
		}
	}

	return nil, krylib.ErrNotImplemented
} // func (pk *PkgApt) Search(string) ([]Package, error)

func (pk *PkgApt) Install(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Install(args ...string) error

func (pk *PkgApt) Remove(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Remove(args ...string) error

func (pk *PkgApt) Update() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Update() error

func (pk *PkgApt) Upgrade() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Upgrade() error

func (pk *PkgApt) ListInstalled() ([]Package, error) {
	return nil, krylib.ErrNotImplemented
} // func (pkg *PkgApt) ListInstalled() ([]Package, error)

func (pk *PkgApt) Clean() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Clean() error

func (pkg *PkgApt) LastUpdate() (time.Time, error) {
	return time.Unix(0, 0), krylib.ErrNotImplemented
} // func (pkg *PkgApt) LastUpdate() (time.Time, error)
