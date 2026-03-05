// 主题切换功能
function initThemeToggle() {
    const themeToggle = document.getElementById('themeToggle');
    const themeIcon = document.getElementById('themeIcon');
    const themeText = document.getElementById('themeText');

    if (!themeToggle) return;

    // 从 localStorage 读取主题偏好
    let isDark = localStorage.getItem('theme') === 'dark';

    // 如果没有保存的偏好，检查当前主题
    if (localStorage.getItem('theme') === null) {
        // 检查是否已经有初始主题设置
        const html = document.documentElement;
        isDark = !html.classList.contains('light-theme-initial') && !html.classList.contains('light-theme');
    }

    // 更新 UI
    function updateThemeUI() {
        const html = document.documentElement;

        if (isDark) {
            html.classList.remove('light-theme', 'light-theme-initial');
            themeIcon.className = 'fas fa-sun';
            themeText.textContent = 'Light';
        } else {
            html.classList.add('light-theme');
            html.classList.remove('light-theme-initial');
            themeIcon.className = 'fas fa-moon';
            themeText.textContent = 'Dark';
        }
    }

    updateThemeUI();

    // 切换主题
    themeToggle.addEventListener('click', function () {
        isDark = !isDark;
        localStorage.setItem('theme', isDark ? 'dark' : 'light');
        updateThemeUI();
    });

    // 监听系统主题变化（仅当没有用户偏好时）
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)');
    prefersDark.addEventListener('change', function (e) {
        if (localStorage.getItem('theme') === null) {
            isDark = e.matches;
            updateThemeUI();
        }
    });
}

// 立即设置主题，防止闪屏（在 DOM 加载前执行）
(function initThemeImmediately() {
    try {
        const savedTheme = localStorage.getItem('theme');
        let shouldBeLight = false;

        if (savedTheme !== null) {
            shouldBeLight = savedTheme === 'light';
        } else {
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)');
            shouldBeLight = !prefersDark.matches;
        }

        if (shouldBeLight) {
            document.documentElement.classList.add('light-theme-initial');
        }
        document.documentElement.classList.add('theme-initialized');
    } catch (e) {
        // localStorage 可能不可用，忽略错误
    }
})();

// 解析文件大小为字节数
function parseSize(sizeStr) {
    if (sizeStr === '-') return -1;
    const units = { 'B': 1, 'KiB': 1024, 'MiB': 1024 * 1024, 'GiB': 1024 * 1024 * 1024, 'TiB': 1024 * 1024 * 1024 * 1024 };
    const match = sizeStr.match(/^([\d.]+)\s*([A-Za-z]+)/);
    if (!match) return parseFloat(sizeStr);
    return parseFloat(match[1]) * (units[match[2]] || 1);
}

// 初始化排序功能
function initSort() {
    const titles = document.querySelectorAll('.th-title');
    titles.forEach(function(title) {
        title.addEventListener('click', function(e) {
            e.stopPropagation();
            const header = this.closest('th');
            const sortBy = header.dataset.sortBy;
            const table = this.closest('.file-table');
            const tbody = table.querySelector('tbody');
            const rows = Array.from(tbody.querySelectorAll('tr'));

            if (!header || !tbody) return;

            const sortState = { by: sortBy, order: 'asc' };

            // 检查当前排序状态
            if (header.classList.contains('sorted')) {
                const currentIcon = header.querySelector('.sort-icon');
                if (currentIcon && currentIcon.classList.contains('fa-sort-up')) {
                    sortState.order = 'desc';
                }
            }

            // 更新排序状态
            const headers = table.querySelectorAll('th.sortable');
            headers.forEach(function(h) {
                h.classList.remove('sorted');
                const icon = h.querySelector('.sort-icon');
                if (icon) {
                    icon.className = 'sort-icon fas fa-sort';
                }
            });

            header.classList.add('sorted');
            const currentIcon = header.querySelector('.sort-icon');
            if (currentIcon) {
                if (sortState.order === 'asc') {
                    currentIcon.className = 'sort-icon fas fa-sort-up';
                } else {
                    currentIcon.className = 'sort-icon fas fa-sort-down';
                }
            }

            const dataRows = rows.filter(function(row) {
                return !row.classList.contains('parent-dir-row');
            });

            dataRows.sort(function(a, b) {
                let aVal, bVal;
                if (sortBy === 'name') {
                    aVal = a.querySelector('.file-name').textContent.trim();
                    bVal = b.querySelector('.file-name').textContent.trim();
                    const aIsDir = a.querySelector('.fa-folder') !== null;
                    const bIsDir = b.querySelector('.fa-folder') !== null;
                    if (aIsDir && !bIsDir) return sortState.order === 'asc' ? -1 : 1;
                    if (!aIsDir && bIsDir) return sortState.order === 'asc' ? 1 : -1;
                } else if (sortBy === 'size') {
                    aVal = parseSize(a.querySelector('.file-size').textContent.trim());
                    bVal = parseSize(b.querySelector('.file-size').textContent.trim());
                } else if (sortBy === 'modified') {
                    aVal = new Date(a.querySelector('.file-modified').textContent.trim());
                    bVal = new Date(b.querySelector('.file-modified').textContent.trim());
                }
                if (aVal < bVal) return sortState.order === 'asc' ? -1 : 1;
                if (aVal > bVal) return sortState.order === 'asc' ? 1 : -1;
                return 0;
            });

            dataRows.forEach(function(row) {
                tbody.appendChild(row);
            });
        });
    });
}

// 初始化搜索功能
function initSearch() {
    const searchInput = document.getElementById('searchInput');
    if (!searchInput) return;
    searchInput.addEventListener('input', function() {
        const filter = this.value.toLowerCase();
        const table = this.closest('.file-table');
        const tbody = table.querySelector('tbody');
        const rows = tbody.querySelectorAll('tr');
        rows.forEach(function(row) {
            if (row.classList.contains('parent-dir-row')) {
                row.style.display = '';
                return;
            }
            const fileName = row.querySelector('.file-name').textContent.trim();
            row.style.display = fileName.toLowerCase().indexOf(filter) > -1 ? '' : 'none';
        });
    });
}

// 初始化代码块复制按钮
function initCodeBlocks() {
    const codeBlocks = document.querySelectorAll('.readme-content pre');
    codeBlocks.forEach(function(pre) {
        let button = pre.querySelector('.code-copy-btn');
        if (!button) {
            button = document.createElement('button');
            button.className = 'code-copy-btn';
            button.textContent = 'Copy';
            pre.insertBefore(button, pre.firstChild);
        }
        button.addEventListener('click', function() {
            const code = pre.querySelector('code');
            if (code) {
                navigator.clipboard.writeText(code.textContent.replace(/\n$/, '')).then(function() {
                    button.textContent = 'Copied!';
                    button.classList.add('copied');
                    setTimeout(function() {
                        button.textContent = 'Copy';
                        button.classList.remove('copied');
                    }, 2000);
                }).catch(function() {
                    button.textContent = 'Failed';
                    setTimeout(function() { button.textContent = 'Copy'; }, 2000);
                });
            }
        });
    });
    // 手动触发 prism 高亮
    if (window.Prism) {
        setTimeout(function() { Prism.highlightAll(); }, 100);
    }
}

// SPA 导航功能
function initSPANavigation() {
    const breadcrumbsNav = document.querySelector('.breadcrumbs');
    const statsDiv = document.querySelector('.stats');

    // 处理目录链接点击（包括 .. 和面包屑）
    document.addEventListener('click', function(e) {
        // 检查是否是文件链接或面包屑链接
        const fileLink = e.target.closest('.file-link');
        const breadcrumbLink = e.target.closest('.breadcrumb-link');

        if (!fileLink && !breadcrumbLink) return;

        // 获取链接地址
        let href = fileLink ? fileLink.getAttribute('href') : breadcrumbLink.getAttribute('href');

        // 跳过外部链接
        if (!href || href.startsWith('http://') || href.startsWith('https://')) {
            return;
        }

        e.preventDefault();

        // 使用 fetch 加载新页面
        fetch(href)
            .then(response => response.text())
            .then(html => {
                // 解析新页面内容
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');

                // 更新标题
                const newTitle = doc.querySelector('title').textContent;
                document.title = newTitle;

                // 注意：Server-Name 是固定的服务器名称，不应该随目录变化
                // 只有标签页标题（document.title）会显示当前目录名

                // 更新面包屑
                const newBreadcrumbs = doc.querySelector('.breadcrumbs').innerHTML;
                if (breadcrumbsNav) {
                    breadcrumbsNav.innerHTML = newBreadcrumbs;
                }

                // 更新统计信息
                const newStats = doc.querySelector('.stats').innerHTML;
                if (statsDiv) {
                    statsDiv.innerHTML = newStats;
                }

                // 先从新页面文档中提取需要的节点（在修改 doc 之前）
                const newTable = doc.querySelector('.file-table');
                const newEmptyDiv = doc.querySelector('.empty-directory');
                const newReadme = doc.querySelector('.readme-container');

                // 更新文件表格
                const fileTable = document.querySelector('.file-table');
                if (newTable) {
                    if (fileTable) {
                        fileTable.replaceWith(newTable.cloneNode(true));
                    }
                } else {
                    // 如果是空目录，替换为空目录提示
                    const emptyDiv = document.querySelector('.empty-directory');
                    if (newEmptyDiv && fileTable) {
                        emptyDiv.replaceWith(newEmptyDiv.cloneNode(true));
                    }
                }

                // 更新/移除 README 区域
                const oldReadme = document.querySelector('.readme-container');
                if (newReadme) {
                    if (oldReadme) {
                        oldReadme.replaceWith(newReadme.cloneNode(true));
                    } else {
                        // 获取新的文件列表位置（如果刚被替换）
                        const currentFileTable = document.querySelector('.file-table');
                        if (currentFileTable) {
                            // 当前页面没有 README，把新的 README 插入到文件列表后面
                            currentFileTable.insertAdjacentElement('afterend', newReadme.cloneNode(true));
                        } else if (statsDiv) {
                            // 空目录时，插入到统计信息后面
                            statsDiv.insertAdjacentElement('afterend', newReadme.cloneNode(true));
                        }
                    }
                } else if (oldReadme) {
                    oldReadme.remove();
                }

                // 更新浏览器历史
                history.pushState({ path: href }, newTitle, href);

                // 重新绑定事件
                rebindEvents();

                // 滚动到顶部
                window.scrollTo(0, 0);
            })
            .catch(err => {
                console.error('Navigation failed:', err);
                // 如果 fetch 失败，回退到传统导航
                window.location.href = href;
            });
    });

    // 处理浏览器后退/前进按钮
    window.addEventListener('popstate', function(e) {
        if (e.state && e.state.path) {
            // 回退到整页加载
            location.reload();
        }
    });
}

// 重新绑定 SPA 导航后的事件
function rebindEvents() {
    initSort();
    initSearch();
    initCodeBlocks();
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    initThemeToggle();
    initSort();
    initSearch();
    initCodeBlocks();
    initSPANavigation();
});
