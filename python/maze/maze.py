from langchain.agents import initialize_agent, AgentType
from langchain.tools import Tool
# 使用更新的导入方式
from langchain_openai import ChatOpenAI
import os
from typing import Tuple

matrix = []
matrix.append(['#', '#', '#', '#', '#', '#', '#'])
matrix.append(['#', 'X', '0', '#', '#', '#', '#'])
matrix.append(['#', '#', '0', '0', '0', '#', '#'])
matrix.append(['#', '#', '#', '#', '0', '#', '#'])
matrix.append(['#', '#', '0', '0', '0', '#', '#'])
matrix.append(['#', '#', '#', '#', '0', '0', '#'])
matrix.append(['#', '#', '#', '#', '#', '#', '#'])

position_x = 1
position_y = 1

# 1. 定义工具函数及其描述
def move_action(tuple_arg: Tuple[int, int]) -> bool:
    """
        player从当前位置尝试移动到新位置
        新位置为 (position_x + tuple_arg[0], position_y + tuple_arg[1])
        每次移动的动作只能是下面4种动作之一:
        (1, 0), (0, 1), (-1, 0), (0, -1)
        :param tuple_arg: 动作
        :type tuple_arg: 元组
        :return: 是否成功移动到新位置
        :rtype: bool
    """
    print(type(tuple_arg))
    my_tuple = eval(tuple_arg)
    print(type(my_tuple))  # 输出: <class 'tuple'>
    move_x = int(my_tuple[0])
    move_y = int(my_tuple[1])
    print(f"action: ({move_x}, {move_y})")

    global position_x
    global position_y
    if move_x * move_y != 0:
        return False
    new_x = position_x + move_x
    new_y = position_y + move_y
    if matrix[new_x][new_y] == '0':
        print(f"从({position_x}, {position_y}) 移动到了 ({new_x}, {new_y})")
        matrix[new_x][new_y] = 'X'
        matrix[position_x][position_y] = '0'
        position_x = new_x
        position_y = new_y
        return True

    return False

move_tool = Tool(
    name="move_action",
    func=move_action,
    description="移动位置，参数为move_x和move_y，例如move_action(1, 0)表示向右移动一步"
)

def observe(*args, **kwargs) -> list:
    """
        从当前位置能观察到的矩阵的一部分, 'X'是你的位置，'#' 都是墙体，只有 '0' 是通道， '-' 表示未知
        :return: 返回 7*7的二维矩阵
        :rtype: list
    """
    views = []
    for x in range(7):
        views.append(['-', '-', '-', '-', '-', '-', '-'])
        for y in range(7):
            if (x >= position_x -1) and (x <= position_x + 1) and (y >= position_y -1) and (y <= position_y + 1):
                views[x][y] = matrix[x][y]

    views[position_x][position_y] = 'X'
    return views

observe_tool = Tool(
    name="observe",
    func=observe,
    description="从player所在位置能观察到的矩阵的情况，无需参数"
)



# 初始化 Kimi 模型
llm = ChatOpenAI(
    model_name="gpt-3.5-turbo", # 或其他 Kimi 模型名
    # openai_api_base=kimi_api_base,
    temperature=0.1
)

agent = initialize_agent(
    tools=[move_tool, observe_tool],  # 工具列表
    llm=llm,
    agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,  # 指定智能体类型
    max_iterations = 200,
    verbose=True  #  verbose=True 可以打印出详细的推理过程

)

# 3. 运行智能体 - 使用新的invoke方法
question = """
你作为player当前处于一个 7*7的 二维矩阵中，'X'是你的位置，'#' 都是墙体，只有 '0' 是通道， '-' 表示未知
初始坐标为(1，1)，你需要设法到达(5，5)
注意: 只能沿着通道移动
你每次可以移动一步，只能是下面4种动作之一
(1, 0), (0, 1), (-1, 0), (0, -1)
请你给出一个完整的移动轨迹，如果可以尝试推测一下矩阵的情况
"""
response = agent.invoke({"input": question})
print(response)