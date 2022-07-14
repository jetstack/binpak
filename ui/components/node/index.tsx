import React, { useState } from "react";
import styles from "./styles.module.scss";

type NodeProps = {
  workload?: any;
  limit?: number;
  hoveredWorkload?: number;
  setHoveredWorkload?: any;
};

const totalBoxes = 200;
const percentagePerBox = 100 / totalBoxes;

export function Node({
  workload = temp,
  limit = 1000,
  hoveredWorkload,
  setHoveredWorkload,
}: NodeProps) {
  function generateBoxes() {
    const list = workload.map(({ value, id }: any, index: number) => ({
      id,
      numberOfBoxes: Math.round(percentage(value, limit) / percentagePerBox),
      color: colorList[index],
    }));

    const componentList = [] as any[];

    list.forEach(({ numberOfBoxes, color, id }: any) => {
      const boxes = new Array(numberOfBoxes).fill(1);
      boxes.forEach((_: any, index: number) => {
        componentList.push(
          <div
            data-id={id}
            key={`${id}-${index}`}
            className={`${styles.pixel} ${
              hoveredWorkload === id ? styles.hovered : ""
            }`}
            style={{
              backgroundColor: color,
            }}
            onMouseOver={() => setHoveredWorkload(id)}
            onMouseOut={() => setHoveredWorkload(undefined)}
          ></div>
        );
      });
    });

    if (componentList.length > totalBoxes) {
      return componentList.map((c, index) => {
        if (index >= totalBoxes) {
          return {
            ...c,
            props: { ...c.props, style: { ...c.props.style, opacity: 0.6 } },
          };
        } else {
          return c;
        }
      });
    } else {
      return [
        ...componentList,
        new Array(totalBoxes - componentList.length)
          .fill(1)
          .map((_: any, index: number) => {
            return <div key={index + "a"} className={`${styles.pixel}`}></div>;
          }),
      ];
    }
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.border}></div>
      <div className={styles.container}>{generateBoxes().map((c) => c)}</div>
    </div>
  );
}
function percentage(partialValue: number, totalValue: number) {
  return (100 * partialValue) / totalValue;
}

const colorList = [
  "#74A8C6",
  "#FFD000",
  "#ED565A",
  "#00FF00",
  "#A742CA",
  "#96C674",
  "#CA427B",
];

const temp = [
  {
    id: 1,
    value: 123,
  },
  {
    id: 2,
    value: 280,
  },
  {
    id: 3,
    value: 280,
  },
  {
    id: 4,
    value: 280,
  },
  {
    id: 5,
    value: 280,
  },
];
