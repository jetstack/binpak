import React, { Fragment, useState } from "react";
import styles from "./styles.module.scss";
import { Tooltip, Fade } from "@mui/material";

type NodeProps = {
  workloads?: any;
  limit?: number;
  hoveredWorkload?: number;
  setHoveredWorkload?: any;
  isCPU?: boolean;
  isAllocated?: boolean;
};

const totalBoxes = 200;
const percentagePerBox = 100 / totalBoxes;

const colorList = [
  "#74A8C6",
  "#3300FF",
  "#FFD000",
  "#ED565A",
  "#00FF00",
  "#A742CA",
  "#96C674",
  "#CA427B",
  "#F887FF",
  "#AFFFFA",
  "#FFEFD0",
];

export function Node({
  workloads,
  limit = 1000,
  isCPU = false,
  isAllocated = false,
  hoveredWorkload,
  setHoveredWorkload,
}: NodeProps) {
  function generateBoxes() {
    //Modify the list of workloads to be used to generate the little percentage boxes
    const list = workloads
      .map((workload: any, index: number) => {
        //Get the value depending on type of resource
        function getValue() {
          let value = 0;
          if (isCPU) {
            value = isAllocated ? workload.limits.cpuM : workload.requests.cpuM;
          } else {
            value = isAllocated
              ? workload.limits?.memoryMi
              : workload.requests?.memoryMi;
          }

          return value;
        }

        const value = getValue();
        const numberOfBoxes =
          Math.round(percentage(value, limit) / percentagePerBox) || 0;

        if (numberOfBoxes > 0) {
          return {
            id: workload.name,
            numberOfBoxes,
            value,
            color: colorList[index],
          };
        }
        return;
      })
      .filter(Boolean);

    const componentList = [] as any[];

    //Generate the percentage boxes if its filled or not
    list.forEach(({ numberOfBoxes, color, value, id }: any) => {
      const boxes = new Array(numberOfBoxes).fill(1);
      boxes.forEach((_: any, index: number) => {
        componentList.push(
          <Tooltip
            title={
              <TooltipContents
                id={id}
                limit={limit}
                value={value}
                isCPU={isCPU}
              />
            }
            arrow
            followCursor
            TransitionComponent={Fade}
            TransitionProps={{ timeout: 0 }}
            key={id + index}
          >
            <div
              className={`${styles.pixel} ${
                hoveredWorkload === id ? styles.hovered : ""
              }`}
              style={{
                backgroundColor: color,
              }}
              onMouseOver={() => setHoveredWorkload(id)}
              onMouseOut={() => {
                setHoveredWorkload(undefined);
              }}
            ></div>
          </Tooltip>
        );
      });
    });

    if (componentList.length > totalBoxes) {
      //If the number of boxes is greater than the total number of boxes,
      //we need to make the last boxes transparent
      return componentList.map((c, index) => {
        if (index >= totalBoxes) {
          return {
            ...c,
            props: {
              ...c.props,
              children: {
                ...c.props.children,
                props: {
                  ...c.props.children.props,
                  style: { ...c.props.children.props.style, opacity: 0.6 },
                },
              },
            },
          };
        } else {
          return c;
        }
      });
    } else {
      //If the list is smaller than the total number of boxes, fill the rest with empty boxes
      const fillBoxes = new Array(totalBoxes - componentList.length)
        .fill(1)
        .map((_: any, index: number) => {
          return <div key={index + "a"} className={styles.pixel}></div>;
        });

      return [...componentList, ...fillBoxes];
    }
  }

  return (
    <div>
      <div className={styles.title}>
        Size: {limit}
        {isCPU ? "m" : "Mi"}
      </div>
      <div className={styles.wrapper}>
        <div className={styles.border}></div>
        <div className={styles.container}>{generateBoxes().map((c) => c)}</div>
      </div>
    </div>
  );
}
function percentage(partialValue: number, totalValue: number) {
  return (100 * partialValue) / totalValue;
}

function TooltipContents({ id, limit, value, isCPU }: any) {
  return (
    <div>
      <h1 className={styles.tooltip}>{id}</h1>
      <p className={styles.tooltipParagraph}>
        <b>Size:</b> {value}
        {isCPU ? "m" : "Mi"}
      </p>
      <p className={styles.tooltipParagraph}>
        <b>Usage:</b> {Math.round(percentage(value, limit) * 100) / 100}%
      </p>
    </div>
  );
}
