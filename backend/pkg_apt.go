// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_apt.go
// -*- mode: go; coding: utf-8; -*-
// Created on 21. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-24 15:48:46 krylon>

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
	cmdApt = "/usr/bin/apt" // nolint: unused
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

*/

var patSearchApt = regexp.MustCompile(`(?mi)^(\S+)\s+-\s+([^\n]+)$`)

func (pk *PkgApt) Search(query string) ([]Package, error) {
	const cmdSearch = "/usr/bin/apt-cache"
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

	var matches = patSearchApt.FindAllStringSubmatch(bufOut.String(), -1)

	if len(matches) == 0 {
		return nil, nil
	}

	var pkList = make([]Package, len(matches))

	for i, m := range matches {
		pkList[i] = Package{
			Name:        m[1],
			Description: m[2],
		}
	}

	return pkList, nil
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
