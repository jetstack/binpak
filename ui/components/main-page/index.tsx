import React, { useState } from "react";
import { Toggle, Node } from "../";
import styles from "./styles.module.scss";

export function MainPage() {
  const [isTypeOfResourceToggled, setIsTypeOfResourceToggled] = useState(false);
  const [isRequested, setIsRequested] = useState(false);
  const [hoveredWorkload, setHoveredWorkload] = useState(undefined);

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
          rightSideItem="Allocated"
          isToggled={isRequested}
          onToggle={setIsRequested}
        />
      </div>
      <div className={styles.list}>
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
        <Node
          hoveredWorkload={hoveredWorkload}
          setHoveredWorkload={setHoveredWorkload}
        />
      </div>
    </div>
  );
}
