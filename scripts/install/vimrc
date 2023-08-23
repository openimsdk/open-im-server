" learn: https://github.com/cubxxw/awesome-cs-course/tree/master/linux
" Read: https://github.com/cubxxw/awesome-cs-course/blob/master/linux/markdown/my_vim.md

"vim的配置关于鼠标滚动滑动"""
if has("autocmd")
  au BufReadPost * if line("'\"") > 1 && line("'\"") <= line("$") | exe "normal! g'\"" | endif
  set mouse=a
endif
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"C，C++ 按F5编译运行
map <F5> :call CompileRunGcc()<CR>
func! CompileRunGcc()
    exec "w"
    if &filetype == 'c'
        exec "!g++ % -o %<"
        exec "! ./%<"
    elseif &filetype == 'cpp'
        exec "!g++ % -o %<"
        exec "! ./%<"
    elseif &filetype == 'java' 
        exec "!javac %" 
        exec "!java %<"
    elseif &filetype == 'sh'
        :!./%
    endif
endfunc
"配置命令"
"C,C++的调试
"''''''''''''''''""""""""""""""""""""""""""""""""""'''''''''''''''''''''
map <F8> :call Rungdb()<CR>
func! Rungdb()
    exec "w"
    exec "!g++ % -g -o %<"
    exec "!gdb ./%<"
endfunc
 
""实用设置
" 设置当文件被改动时自动载入
set autoread
" quickfix模式
autocmd FileType c,cpp map <buffer> <leader><space> :w<cr>:make<cr>
"代码补全 
set completeopt=preview,menu

" common configure
"set noswapfile
set mouse=a " 激活鼠标使用
set wrap " 自动换行
set sw=4 " 设置软宽度
set gdefault " 行内替换
set nu
set showmatch " 高亮显示括号匹配    
set expandtab " 使用空格来替换 Tab
set tabstop=4 " 设置 Tab 长度为 4 空格    
set shiftwidth=4 " 设置自动缩进长度为 4 空格    
set autoindent " 继承前一行的缩进方式，适用于多行注释" 
set autowrite " 自动保存
set nocompatible " 关闭 vi 兼容模式
set history=1000 " 设置历史记录步数
set confirm " 在处理未保存或只读文件时，弹出确认

" 搜索逐字符高亮
set hlsearch
set incsearch

" 从不备份
set nobackup
set noswapfile

let g:indentLine_enabled=0

" golang configure ====> start
let g:go_highlight_methods = 1                 
let g:go_highlight_operators = 1    
let g:go_highlight_build_constraints = 1    
let g:go_fmt_autosave = 1
let g:go_version_warning = 1    
let g:go_autodetect_gopath = 1    
let g:go_highlight_types = 1    
let g:go_highlight_fields = 1    
let g:go_highlight_functions = 1    
let g:go_highlight_function_calls = 1                                             
let g:go_highlight_extra_types = 1                                                               
let g:go_highlight_generate_tags = 1    
let g:go_def_mode = 'gopls'
let g:go_gopls_enabled = 1
let g:go_guru_enabled = 1
let g:go_fmt_experimental = 1    
let g:go_def_mapping_enabled = 1    
let g:go_build_tags = '-v'    
let g:go_fmt_command = "goimports"                       
let g:go_list_type = "quickfix"    
let g:go_def_mapping_enable = 1

map <2-LeftMouse> :GoDef <CR>
map <RightMouse> :GoDefPop <CR>
map <C-N> :GoCallers <CR>
map <C-M> :GoCallees <CR>
map <ScrollWheelUp> <C-Y>
map <ScrollWheelDown> <C-E>
unmap <CR>
map <CR> :GoDef <CR>
map <C-I> :GoDefPop <CR>
map <S-K> :GoDoc<cr>
map <S-M> :GoInfo<cr>
map <S-T> :GoDefType<cr>
map <S-L> :GoAddTag<cr>
map <S-P> :GoImplements<cr>
map <S-R> :GoRename<cr>
map <S-F> :GoFillStruct<cr>
map <S-C> :GoCallers<cr>
map <S-H> :GoSameIdsToggle<cr>

augroup go
  autocmd!

  " Show by default 4 spaces for a tab
  autocmd BufNewFile,BufRead *.go setlocal noexpandtab tabstop=4 shiftwidth=4

  " :GoDef but opens in a vertical split
  autocmd FileType go nmap <Leader>v <Plug>(go-def-vertical)
  " :GoDef but opens in a horizontal split
  autocmd FileType go nmap <Leader>s <Plug>(go-def-split)

  " :GoAlternate  commands :A, :AV, :AS and :AT
  autocmd Filetype go command! -bang A call go#alternate#Switch(<bang>0, 'edit')
  autocmd Filetype go command! -bang AV call go#alternate#Switch(<bang>0, 'vsplit')
  autocmd Filetype go command! -bang AS call go#alternate#Switch(<bang>0, 'split')
  autocmd Filetype go command! -bang AT call go#alternate#Switch(<bang>0, 'tabe')
augroup END

" build_go_files is a custom function that builds or compiles the test file.
" It calls :GoBuild if its a Go file, or :GoTestCompile if it's a test file
function! s:build_go_files()
  let l:file = expand('%')
  if l:file =~# '^\f\+_test\.go$'
    call go#test#Test(0, 1)
  elseif l:file =~# '^\f\+\.go$'
    call go#cmd#Build(0)
  endif
endfunction
" golang configure ====> end

:inoremap ( ()<ESC>i
:inoremap ) <c-r>=ClosePair(')')<CR>
:inoremap { {<CR>}<ESC>O
:inoremap } <c-r>=ClosePair('}')<CR>
:inoremap [ []<ESC>i
:inoremap ] <c-r>=ClosePair(']')<CR>
:inoremap " ""<ESC>i
:inoremap ' ''<ESC>i
function! ClosePair(char)
    if getline('.')[col('.') - 1] == a:char
        return "\<Right>"
    else
        return a:char
    endif
endfunction
filetype plugin indent on
"打开文件类型检测, 加了这句才可以用智能补全

set completeopt=longest,menu

"""""""""""""""""""""""""""""""""""""""""""""""
  noremap <F6> :set nu
  noremap <F4> :!python3 a  
  noremap <F7> :set ai
  noremap <F8> :syntax on
  set cursorline     " 突出显示当前行
  set magic          " 设置魔术
 "" noremap <F5> :! g++ -o a 
nnoremap <C-J> <C-W><C-J>
nnoremap <C-K> <C-W><C-K>
nnoremap <C-L> <C-W><C-L>
nnoremap <C-H> <C-W><C-H>
nmap wj :resize -3<CR>
nmap wk :resize +3<CR>
nmap wh :vertical resize -3<CR>
nmap wl :vertical resize +3<CR>
set guifont=Droid\ Sans\ Mono\ Nerd\ Font\ Complete:h18 " 设置字体
set guicursor=n-v-c:ver5    " 设置光标为竖线
set number ""# 显示行号
set autowrite "" # 自动保存
set ruler  ""# 显示打开状态栏标尺
set cursorline "" # 突出显示当前行

set showmatch "" # 匹配光标所经过的括号等.
set showcmd ""# 命令行显示输入的命令
set showmode ""命令行显示vim当前模式
set showtabline=0           " 隐藏Tab栏
set laststatus=2 """"'vim 窗口底部显示永久状态栏，显示文件名，行号，列号等.
let mapleader = ","      " 定义<leader>键
set nocompatible         " 设置不兼容原始vi模式
filetype on              " 设置开启文件类型侦测
filetype plugin on       " 设置加载对应文件类型的插件
set noeb                 " 关闭错误的提示
syntax enable            " 开启语法高亮功能
syntax on                " 自动语法高亮
set cmdheight=2          " 设置命令行的高度
set showcmd              " select模式下显示选中的行数
set ruler                " 总是显示光标位置
set laststatus=2         " 总是显示状态栏
set number               " 开启行号显示
set cursorline           " 高亮显示当前行
set whichwrap+=<,>,h,l   " 设置光标键跨行
set ttimeoutlen=0        " 设置<ESC>键响应时间
set virtualedit=block,onemore   " 允许光标出现在最后一个字符的后面

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" 代码缩进和排版
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
set autoindent           " 设置自动缩进
set cindent              " 设置使用C/C++语言的自动缩进方式
set cinoptions=g0,:0,N-s,(0    " 设置C/C++语言的具体缩进方式
set smartindent          " 智能的选择对其方式
filetype indent on       " 自适应不同语言的智能缩进
set expandtab            " 将制表符扩展为空格
set tabstop=4            " 设置编辑时制表符占用空格数
set shiftwidth=4         " 设置格式化时制表符占用空格数
set softtabstop=4        " 设置4个空格为制表符
set smarttab             " 在行和段开始处使用制表符
set nowrap               " 禁止折行
set backspace=2          " 使用回车键正常处理indent,eol,start等
set sidescroll=10        " 设置向右滚动字符数
set nofoldenable         " 禁用折叠代码

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" 代码补全
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
set wildmenu             " vim自身命名行模式智能补全
set completeopt-=preview " 补全时不显示窗口，只显示补全列表

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" 搜索设置
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
set hlsearch            " 高亮显示搜索结果
set incsearch           " 开启实时搜索功能
set ignorecase          " 搜索时大小写不敏感

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" 缓存设置
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
set nobackup            " 设置不备份
set noswapfile          " 禁止生成临时文件
set autoread            " 文件在vim之外修改过，自动重新读入
set autowrite           " 设置自动保存
set confirm             " 在处理未保存或只读文件的时候，弹出确认

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"调整窗
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
nmap wj :resize -3<CR>
nmap wk :resize +3<CR>
nmap wh :vertical resize -3<CR>
nmap wl :vertical resize +3<CR>

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"编码"
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
set langmenu=zh_CN.UTF-8
set helplang=cn
set termencoding=utf-8
set encoding=utf8
set fileencodings=utf8,ucs-bom,gbk,cp936,gb2312,gb18030

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"代码补全“
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
inoremap ' ''<ESC>i
inoremap " ""<ESC>i
inoremap ( ()<ESC>i
inoremap [ []<ESC>i
inoremap { {<CR>}<ESC>O
 : set nu
 : set ai
 : syntax on
 : filetype on
  set tabstop=4
  set ignorecase
noremap <C-h> <C -w>h
noremap <C-j> <C -w>j
noremap <C-k> <C -w>r
noremap <C-l> <C -w>l
" Specify a directory for plugins
" - For Neovim: stdpath('data') . '/plugged'
" - Avoid using standard Vim directory names like 'plugin'
let g:coc_disable_startup_warning = 1 

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"tmux"
"复用终端、分屏"
"let g:EasyMotion_startofline = 0 " keep cursor colum when JK motion
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
map <Leader><leader>h <Plug>(easymotion-linebackward)
map <Leader><Leader>j <Plug>(easymotion-j)
map <Leader><Leader>k <Plug>(easymotion-k)
map <Leader><leader>l <Plug>(easymotion-lineforward)
" 重复上一次操作, 类似repeat插件, 很强大
map <Leader><leader>. <Plug>(easymotion-repeat)
nmap s <Plug>(easymotion-s)
