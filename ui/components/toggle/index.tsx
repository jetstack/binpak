import React from "react";
import styles from "./styles.module.scss";

type ToggleProps = {
  leftSideItem: string;
  rightSideItem: string;
  isToggled: boolean;
  onToggle: (isOn: any) => void;
};

export function Toggle({
  leftSideItem,
  rightSideItem,
  isToggled,
  onToggle,
}: ToggleProps) {
  return (
    <div
      className={`${styles.wrapper} ${
        isToggled ? styles.isToggled : styles.isNotToggled
      }`}
      onClick={() => onToggle((state: boolean) => !state)}
    >
      <div className={styles.leftItem}>{leftSideItem}</div>
      <div className={styles.toggleContainer}>
        <div className={styles.toggle}></div>
      </div>
      <div className={styles.rightItem}>{rightSideItem}</div>
    </div>
  );
}
