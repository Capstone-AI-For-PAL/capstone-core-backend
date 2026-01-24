import requests
import sys

# Configuration
URL = "http://localhost:5000/generate"
OUTPUT_FILE = "test_presentation.pdf"

# Sample Marp Markdown
markdown_payload = """---
marp: true
theme: gaia
class: lead
backgroundColor: #f0f4f8
paginate: true
style: |
  section {
    font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  }
  h1 {
    color: #2c3e50;
    margin-top: 0;
  }
  h2 {
    color: #3498db;
    margin-top: 0;
  }
---

![bg left:40%](https://images.unsplash.com/photo-1544367563-12123d8965cd?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80)

# **Unlocking Longevity**

### Science, Strategy, and Living Better Longer

<br>

**Overview**
* The Science of Aging
* Actionable Strategies
* The Blue Zones

---

## **The Science: Why Do We Age?**

Current research identifies key "Hallmarks of Aging" that drive cellular decline.

* **Telomere Attrition**
    * Protective caps on our chromosomes shorten as cells divide, eventually stopping replication.
* **Mitochondrial Dysfunction**
    * The "power plants" of our cells lose efficiency over time, reducing energy and increasing oxidative stress.
* **Cellular Senescence**
    * "Zombie cells" stop dividing but don't die; they accumulate and damage surrounding healthy tissue.

> **Key Insight:** Aging is now viewed by many researchers not as an inevitable decline, but as a treatable biological condition.

---

## **Blueprints for a Longer Life**

Based on data from the **Blue Zones** (regions with the most centenarians).

| Category | Strategy |
| :--- | :--- |
| **Nutrition** | **Plant-Slant Diet:** High intake of beans, nuts, and greens. Stop eating when 80% full (*Hara Hachi Bu*). |
| **Movement** | **Natural Movement:** Gardening, walking, and manual tasks rather than intense gym sessions. |
| **Recovery** | **Sleep & Stress:** 7-9 hours of quality sleep. Regular downshifting (meditation, naps, prayer). |
| **Connection** | **Social Circles:** Strong ties to family and friends significantly boost life expectancy. |

<br>

**The Goal:** Increase *Healthspan* (years spent in good health), not just *Lifespan*.
"""

def run_test(): 
    print(f"🚀 Sending request to {URL}...")
    try:
        response = requests.post(
            URL, 
            json={"markdown": markdown_payload},
            timeout=30  # Wait up to 30 seconds for PDF generation
        )
        
        # Check if the request was successful (HTTP 200)
        if response.status_code == 200:
            with open(OUTPUT_FILE, "wb") as f:
                f.write(response.content)
            print(f"✅ Success! PDF saved to: {OUTPUT_FILE}")
            print("   Please open the file to verify the content.")
        else:
            print(f"❌ Failed with Status Code: {response.status_code}")
            print(f"   Response: {response.text}")

    except requests.exceptions.ConnectionError:
        print("❌ Connection Error: Could not connect to the service.")
        print("   Is the Docker container running? Check with: docker ps")
    except Exception as e:
        print(f"❌ An error occurred: {str(e)}")

if __name__ == "__main__": 
    run_test()
