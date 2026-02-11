" Basic settings
set number                      " Show line numbers
set relativenumber              " Show relative line numbers
set tabstop=2                   " Number of spaces tabs count for
set shiftwidth=2                " Number of spaces to use for autoindent
set expandtab                   " Use spaces instead of tabs
set autoindent                  " Auto-indent new lines
set smartindent                 " Smart auto-indenting
syntax on                       " Enable syntax highlighting

" Search settings
set hlsearch                    " Highlight search results
set incsearch                   " Incremental search
set ignorecase                  " Case insensitive search
set smartcase                   " Case sensitive if uppercase present

" UI settings
set showmatch                   " Show matching brackets
set ruler                       " Show cursor position
set wildmenu                    " Visual autocomplete for command menu
set laststatus=2                " Always show status line
set showcmd                     " Show command in bottom bar

" Performance
set lazyredraw                  " Don't redraw while executing macros

" Enable mouse support
set mouse=a

" Backspace behavior
set backspace=indent,eol,start

" File type detection
filetype plugin indent on

" Color scheme (if available)
try
  colorscheme desert
catch
endtry
