import React, { useState } from "react";
import { Toggle, Node } from "../";
import styles from "./styles.module.scss";
import { mockResponse } from "./constants";

export function MainPage() {
  const [isTypeOfResourceToggled, setIsTypeOfResourceToggled] = useState(false);
  const [isAllocated, setIsAllocated] = useState(false);
  const [hoveredWorkload, setHoveredWorkload] = useState(undefined);

  const response = mockResponse;
  const isCPU = isTypeOfResourceToggled;

  return (
    <div className={styles.container}>
      <div className={styles.logo}>
        <img src="./images/binpak.svg" alt="" />
      </div>
      <div className={styles.content}>
        <p>
          This is an introduction text to explain what binpak is. This is an
          introduction text to explain what binpak is. This is an introduction
          text to explain what binpak is. This is an introduction text to
          explain what binpak is.
        </p>
      </div>

      <div className={styles.togglesContainer}>
        <Toggle
          leftSideItem="Memory"
          rightSideItem="CPU"
          isToggled={isTypeOfResourceToggled}
          onToggle={setIsTypeOfResourceToggled}
        />
        <div className={styles.divider}></div>
        <Toggle
          leftSideItem="Requested"
          rightSideItem="Limits"
          isToggled={isAllocated}
          onToggle={setIsAllocated}
        />
      </div>
      <h3>Cluster name: {response.clusterName}</h3>
      <div className={styles.list}>
        {response.instances.map((instance: any) => (
          <Node
            key={instance.name}
            workloads={instance.workloads}
            limit={isCPU ? instance.capacity.cpuM : instance.capacity.memoryMi}
            isCPU={isCPU}
            isAllocated={isAllocated}
            hoveredWorkload={hoveredWorkload}
            setHoveredWorkload={setHoveredWorkload}
          />
        ))}
      </div>
    </div>
  );
}
