#!/usr/bin/env python3
"""
同步文档目录脚本

功能：
1. 扫描 docs/guide-zh-CN/ 下所有 md 文件
2. 读取每个文件的一级标题作为显示名称
3. 将未列入 docs/guide-zh-CN/README.md 的文件按分类规则自动添加
4. 更新 docs/guide-zh-CN/README.md
5. 同步到 server/resource/public/docs/README.md 和 sidebar.md

用法:
    python3 sync_docs_sidebar.py [project_root]

其中 project_root 是 hotgo 项目的根目录路径，不传则自动查找。
"""

import os
import re
import sys

# =============================================================================
# 文件名前缀 → 所属分类的映射规则
# 键为文件名前缀，值为该文件应归入的分类标题
# =============================================================================
PREFIX_TO_SECTION = {
    "start-":  "介绍安装",
    "sys-":    "系统开发",
    "dev-":    "系统开发",
    "addon-":  "插件模块开发",
    "code-":   "生成代码",
    "web-":    "前端开发",
    "append-": "附录",
}


def find_project_root(start_path):
    """从给定路径向上查找项目根目录"""
    current = os.path.abspath(start_path)
    while current != os.path.dirname(current):
        if os.path.isfile(os.path.join(current, "docs", "guide-zh-CN", "README.md")):
            return current
        current = os.path.dirname(current)
    return None


def extract_title(filepath):
    """从 md 文件中提取第一个标题作为显示名称。

    依次尝试：
    1. 第一个 # 标题
    2. 第一个 ## 标题
    3. 文件名（去掉扩展名）
    """
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                m = re.match(r'^#{1,2}\s+(.+)$', line)
                if m:
                    return m.group(1).strip()
    except Exception:
        pass
    # 回退：用文件名
    name = os.path.basename(filepath)
    return name.rsplit(".", 1)[0]


def classify_file(filename):
    """根据文件名前缀判断所属分类，返回分类标题或 None"""
    for prefix, section in PREFIX_TO_SECTION.items():
        if filename.startswith(prefix):
            return section
    return None


def parse_source_readme(content):
    """解析源 README.md，提取分类和链接条目。

    返回列表：[(heading_level, heading_text, [(link_text, link_target), ...])]
    """
    sections = []
    current_heading = None
    current_level = None
    current_items = []

    for line in content.splitlines():
        heading_match = re.match(r'^(#{2,4})\s+(.+)$', line)
        if heading_match:
            if current_heading is not None:
                sections.append((current_level, current_heading, current_items))
            current_level = len(heading_match.group(1))
            current_heading = heading_match.group(2).strip()
            current_items = []
            continue

        link_match = re.match(r'^-\s+\[(.+?)\]\((.+?)\)\s*$', line)
        if link_match:
            current_items.append((link_match.group(1), link_match.group(2)))

    if current_heading is not None:
        sections.append((current_level, current_heading, current_items))

    return sections


def get_listed_files(sections):
    """获取 README 中已列出的所有本地 md 文件名集合"""
    listed = set()
    for _, _, items in sections:
        for _, target in items:
            if ".." not in target:
                listed.add(target)
    listed.add("README.md")
    return listed


def discover_unlisted_files(guide_dir, listed_files):
    """发现目录中存在但未被列出的 md 文件"""
    all_md = set()
    for f in os.listdir(guide_dir):
        if f.endswith(".md"):
            all_md.add(f)
    return sorted(all_md - listed_files)


def add_unlisted_to_sections(sections, unlisted_files, guide_dir):
    """将未列出的文件按分类规则插入到对应分类中。

    返回 (updated_sections, unclassified_files)
    """
    # 建立分类标题 → sections 索引的映射
    section_map = {}
    for i, (level, heading, items) in enumerate(sections):
        section_map[heading] = i

    unclassified = []
    added_count = 0

    for filename in unlisted_files:
        section_name = classify_file(filename)
        if section_name is None or section_name not in section_map:
            unclassified.append(filename)
            continue

        title = extract_title(os.path.join(guide_dir, filename))
        idx = section_map[section_name]
        level, heading, items = sections[idx]
        items.append((title, filename))
        sections[idx] = (level, heading, items)
        added_count += 1

    if added_count > 0:
        print(f"📝 已将 {added_count} 个新文件添加到目录中")

    return sections, unclassified


def generate_source_readme(sections):
    """生成 docs/guide-zh-CN/README.md 内容"""
    lines = ["## 目录", ""]

    for level, heading, items in sections:
        if heading == "目录":
            continue

        prefix = "#" * level
        lines.append(f"{prefix} {heading}")
        lines.append("")
        for link_text, link_target in items:
            lines.append(f"- [{link_text}]({link_target})")
        lines.append("")
        lines.append("")

    result = "\n".join(lines).rstrip("\n") + "\n"
    return result


def generate_docs_readme(sections):
    """生成 server/resource/public/docs/README.md 内容"""
    lines = ["## 目录", ""]

    for level, heading, items in sections:
        if heading == "目录":
            continue

        prefix = "#" * level
        lines.append(f"{prefix} {heading}")
        lines.append("")
        for link_text, link_target in items:
            if ".." in link_target:
                continue
            lines.append(f"- [{link_text}](guide-zh-CN/{link_target})")
        lines.append("")
        lines.append("")

    result = "\n".join(lines).rstrip("\n") + "\n"
    return result


def generate_sidebar(sections):
    """生成 server/resource/public/docs/sidebar.md 内容"""
    lines = []
    first_section = True

    for level, heading, items in sections:
        if heading == "目录":
            continue

        if not first_section:
            lines.append("")
            lines.append("")
        first_section = False

        lines.append(f"- {heading}")
        lines.append("")
        for link_text, link_target in items:
            if ".." in link_target:
                lines.append(f"    - [{link_text}](README.md)")
            else:
                lines.append(f"    - [{link_text}](guide-zh-CN/{link_target})")

    lines.append("")
    return "\n".join(lines)


def validate_links(sections, guide_dir):
    """验证所有链接目标文件是否存在"""
    missing = []
    for _, _, items in sections:
        for link_text, link_target in items:
            if ".." in link_target:
                continue
            if not os.path.isfile(os.path.join(guide_dir, link_target)):
                missing.append((link_text, link_target))
    return missing


def main():
    if len(sys.argv) < 2:
        project_root = find_project_root(os.getcwd())
        if project_root is None:
            print("错误: 无法自动找到项目根目录，请提供项目根目录路径作为参数。", file=sys.stderr)
            print(f"用法: python3 {sys.argv[0]} <project_root>", file=sys.stderr)
            sys.exit(1)
    else:
        project_root = os.path.abspath(sys.argv[1])

    guide_dir = os.path.join(project_root, "docs", "guide-zh-CN")
    source_readme_path = os.path.join(guide_dir, "README.md")
    target_readme = os.path.join(project_root, "server", "resource", "public", "docs", "README.md")
    target_sidebar = os.path.join(project_root, "server", "resource", "public", "docs", "sidebar.md")

    if not os.path.isfile(source_readme_path):
        print(f"错误: 源文件不存在: {source_readme_path}", file=sys.stderr)
        sys.exit(1)

    # 1. 读取并解析源 README
    with open(source_readme_path, "r", encoding="utf-8") as f:
        content = f.read()

    sections = parse_source_readme(content)

    # 2. 发现未列出的文件
    listed_files = get_listed_files(sections)
    unlisted = discover_unlisted_files(guide_dir, listed_files)

    if unlisted:
        print(f"🔍 发现 {len(unlisted)} 个未列入目录的文件:")
        for f in unlisted:
            print(f"   - {f}")
        print()

        # 3. 自动分类并添加到 sections
        sections, unclassified = add_unlisted_to_sections(sections, unlisted, guide_dir)

        if unclassified:
            print(f"⚠️  以下文件无法自动分类，请手动添加:")
            for f in unclassified:
                print(f"   - {f}")
            print()

        # 4. 更新源 README
        source_content = generate_source_readme(sections)
        with open(source_readme_path, "w", encoding="utf-8") as f:
            f.write(source_content)
        print(f"✅ 已更新源文件: {source_readme_path}")
    else:
        print("✅ 所有文件均已在目录中列出")

    # 5. 验证链接
    missing = validate_links(sections, guide_dir)
    if missing:
        print("\n⚠️  以下链接目标文件不存在:")
        for link_text, link_target in missing:
            print(f"   - [{link_text}]({link_target})")
        print()

    # 6. 生成并写入 docsify 文件
    os.makedirs(os.path.dirname(target_readme), exist_ok=True)

    readme_content = generate_docs_readme(sections)
    with open(target_readme, "w", encoding="utf-8") as f:
        f.write(readme_content)
    print(f"✅ 已同步: {target_readme}")

    sidebar_content = generate_sidebar(sections)
    with open(target_sidebar, "w", encoding="utf-8") as f:
        f.write(sidebar_content)
    print(f"✅ 已同步: {target_sidebar}")

    if missing:
        sys.exit(1)
    else:
        print("\n✅ 全部完成，所有链接验证通过。")


if __name__ == "__main__":
    main()
