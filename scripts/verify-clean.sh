#!/bin/bash

# BookCommunity - 验证原作者痕迹是否清除

echo "🔍 开始验证 BookCommunity 项目..."
echo ""

# 检查 Git 历史
echo "1️⃣  检查 Git 提交历史:"
COMMIT_COUNT=$(git log --oneline | wc -l)
echo "   提交数量: $COMMIT_COUNT"
if [ "$COMMIT_COUNT" -eq 1 ]; then
    echo "   ✅ 通过：仅有1条初始提交"
else
    echo "   ❌ 失败：提交数量不为1"
fi
echo ""

# 检查提交者信息
echo "2️⃣  检查提交者信息:"
AUTHORS=$(git log --format="%an <%ae>" | sort -u)
echo "   提交者: $AUTHORS"
if echo "$AUTHORS" | grep -q "doraemon"; then
    echo "   ❌ 失败：发现原作者信息"
else
    echo "   ✅ 通过：无原作者痕迹"
fi
echo ""

# 检查模块路径
echo "3️⃣  检查 go.mod 模块路径:"
MODULE_PATH=$(head -1 go.mod | awk '{print $2}')
echo "   模块路径: $MODULE_PATH"
if echo "$MODULE_PATH" | grep -q "Doraemonkeys/douyin2"; then
    echo "   ❌ 失败：仍使用原模块路径"
else
    echo "   ✅ 通过：已使用新模块路径"
fi
echo ""

# 检查代码中的导入路径
echo "4️⃣  检查代码中的导入路径:"
OLD_IMPORTS=$(grep -r "github.com/Doraemonkeys/douyin2" . --include="*.go" | wc -l)
echo "   原路径引用数量: $OLD_IMPORTS"
if [ "$OLD_IMPORTS" -eq 0 ]; then
    echo "   ✅ 通过：无原路径引用"
else
    echo "   ❌ 失败：发现 $OLD_IMPORTS 处原路径引用"
fi
echo ""

# 检查远程仓库
echo "5️⃣  检查远程仓库配置:"
REMOTE=$(git remote -v | head -1)
if [ -z "$REMOTE" ]; then
    echo "   远程仓库: 未配置"
    echo "   ✅ 通过：无原仓库关联"
elif echo "$REMOTE" | grep -q "Doraemonkeys/douyin2"; then
    echo "   远程仓库: $REMOTE"
    echo "   ❌ 失败：仍关联原仓库"
else
    echo "   远程仓库: $REMOTE"
    echo "   ✅ 通过：已关联新仓库"
fi
echo ""

# 检查许可证
echo "6️⃣  检查许可证文件:"
if [ -f "LICENSE" ]; then
    COPYRIGHT=$(grep "Copyright" LICENSE | head -1)
    echo "   版权: $COPYRIGHT"
    if echo "$COPYRIGHT" | grep -q "BookCommunity"; then
        echo "   ✅ 通过：使用新版权声明"
    else
        echo "   ⚠️  警告：版权声明可能需要更新"
    fi
else
    echo "   ❌ 失败：LICENSE 文件不存在"
fi
echo ""

# 总结
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 验证总结"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 计算通过率
PASS_COUNT=0
TOTAL_COUNT=6

[ "$COMMIT_COUNT" -eq 1 ] && PASS_COUNT=$((PASS_COUNT + 1))
! echo "$AUTHORS" | grep -q "doraemon" && PASS_COUNT=$((PASS_COUNT + 1))
! echo "$MODULE_PATH" | grep -q "Doraemonkeys/douyin2" && PASS_COUNT=$((PASS_COUNT + 1))
[ "$OLD_IMPORTS" -eq 0 ] && PASS_COUNT=$((PASS_COUNT + 1))
[ -z "$REMOTE" ] || ! echo "$REMOTE" | grep -q "Doraemonkeys/douyin2" && PASS_COUNT=$((PASS_COUNT + 1))
[ -f "LICENSE" ] && PASS_COUNT=$((PASS_COUNT + 1))

echo "通过项目: $PASS_COUNT / $TOTAL_COUNT"
echo ""

if [ "$PASS_COUNT" -eq "$TOTAL_COUNT" ]; then
    echo "🎉 恭喜！所有检查通过，原作者痕迹已完全清除！"
    echo ""
    echo "下一步："
    echo "1. 修改 Git 用户信息"
    echo "2. 更新 go.mod 模块路径为您的 GitHub 用户名"
    echo "3. 创建 GitHub 仓库并推送"
else
    echo "⚠️  部分检查未通过，请参考上述提示进行修复"
fi

echo ""
